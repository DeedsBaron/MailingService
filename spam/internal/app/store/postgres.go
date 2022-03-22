package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"spam/internal/app/config"
	"spam/internal/app/validation"
	"spam/pkg/utils"
	"strconv"
	"sync"
	"time"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func (st *Postgres) AddUser(ctx context.Context, user *validation.User) error {
	q := `INSERT INTO клиент (phone_num, mobile_code, tag, timezone_abbr)
		VALUES ($1, $2, $3, $4);`

	_, err := st.pool.Exec(ctx, q, user.PhoneNum, user.MobileCode, user.Tag, user.TimezoneAbbrev)
	if err != nil {
		return err
	}
	return nil
}

func (st *Postgres) UpdateMailingList(ctx context.Context, ml *validation.MailingList, ID int) error {
	q := `UPDATE рассылка 
		SET launch_date = COALESCE($1,launch_date),
		message = COALESCE($2,message), 
		filter = COALESCE($3,filter), 
		finish_date = COALESCE($4,finish_date)
		WHERE id = $5;`

	commandTag, err := st.pool.Exec(ctx, q, ml.LaunchDate, ml.Message, ml.Filter, ml.FinishDate, ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("Not valid user ID")
	}
	return nil
}

func (st *Postgres) UpdateUser(ctx context.Context, user *validation.User, ID int) error {
	q := `UPDATE клиент 
		SET phone_num = COALESCE($1,phone_num),
		mobile_code = COALESCE($2,mobile_code), 
		tag = COALESCE($3,tag), 
		timezone_abbr = COALESCE($4,timezone_abbr)
		WHERE id = $5`

	comandTag, err := st.pool.Exec(ctx, q, user.PhoneNum, user.MobileCode, user.Tag, user.TimezoneAbbrev, ID)
	if err != nil {
		return err
	}
	if comandTag.RowsAffected() == 0 {
		return errors.New("Not valid user ID")
	}
	return nil
}

func (st *Postgres) DeleteUser(ctx context.Context, ID int) error {
	q := `DELETE FROM клиент
		WHERE id = $1;`
	comandTag, err := st.pool.Exec(ctx, q, ID)
	if err != nil {
		return err
	}
	if comandTag.RowsAffected() == 0 {
		return errors.New("Not valid user ID")
	}
	return nil
}

func (st *Postgres) DeleteMailingList(ctx context.Context, ID int) error {
	q := `DELETE FROM рассылка
		WHERE id = $1`
	comandTag, err := st.pool.Exec(ctx, q, ID)
	if err != nil {
		return err
	}
	if comandTag.RowsAffected() == 0 {
		return errors.New("Not valid user ID")
	}
	return nil
}

func (st *Postgres) AddMailingList(ctx context.Context, ml *validation.MailingList) error {
	q := `INSERT INTO рассылка (launch_date, message, filter, finish_date)
		VALUES ($1, $2, $3, $4)
		RETURNING "id";`
	err := st.pool.QueryRow(ctx, q, ml.LaunchDate, ml.Message, ml.Filter, ml.FinishDate).Scan(&ml.ID)
	if err != nil {
		return err
	}
	return nil
}

// Body structure of Post method to extern API
type Body struct {
	Id    int
	Phone int
	Text  string
	mu    sync.Mutex
}

func (st *Postgres) StartSpam(ctx context.Context, w http.ResponseWriter, ml *validation.MailingList, conf *config.Config, logger *logrus.Logger) {
	var deadline time.Duration
	//deadline for MailingList
	if ml.LaunchDateTime.Before(time.Now()) && ml.FinishDateTime.After(time.Now()) {
		deadline = ml.FinishDateTime.Sub(time.Now())
		logger.Info("Created a request to send MailingList: ", ml.ID, " to Clients with filter: ", *ml.Filter, " with deadline: "+utils.FmtDuration(deadline))
	} else if ml.LaunchDateTime.After(time.Now()) {
		logger.Info("Created a request to send MailingList:", ml.ID, "to Client with filter:", ml.Filter, "it will start:", ml.LaunchDateTime.Format(time.RFC822))
		//wait until ml.LaunchDateTime.Before(time.Now())
		for ml.LaunchDateTime.After(time.Now()) {
			time.Sleep(1 * time.Second)
		}
		deadline = ml.FinishDateTime.Sub(time.Now())
	} else if ml.FinishDateTime.Before(time.Now()) {
		logger.Error("Cant create request for MailingList:", ml.ID, "to Client with filter:", ml.Filter, " wrong date params!")
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline))
	defer cancel()

	q := `SELECT клиент.id
		FROM рассылка, клиент
		WHERE клиент.tag = split_part(рассылка.filter, ';', 1)
			AND клиент.mobile_code = split_part(рассылка.filter, ';', 2)
			AND рассылка.id = $1;`
	rows, err := st.pool.Query(ctx, q, ml.ID)
	if err != nil {
		utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
		logger.Error(err.Error())
		return
	}

	i := 0
	messageID := 0
	//rows - all clients wich satisfy filter
	for rows.Next() {
		select {
		case <-ctx.Done():
			logger.Info("MailingList:", ml.ID, " has finished due to deadline!")
			return
		default:
			spam := &Body{
				Text: *ml.Message,
			}

			rows.Scan(&spam.Id)

			q = `INSERT INTO сообщение (create_date, status, mailinglist_id, client_id)
			VALUES ($1, $2, $3, $4)
			RETURNING "id";`
			err = st.pool.QueryRow(ctx, q, time.Now().Format(time.RFC822), false, ml.ID, spam.Id).Scan(&messageID)
			if err != nil {
				utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
				logger.Error(err.Error())
				return
			}
			if rows.Err() != nil {
				utils.HttpErrorWithoutBackSlashN(w, rows.Err().Error(), http.StatusBadRequest)
				logger.Error(rows.Err().Error())
				return
			}

			err = st.GetPhoneNum(ctx, spam.Id, &spam.Phone)
			if err != nil {
				utils.HttpErrorWithoutBackSlashN(w, rows.Err().Error(), http.StatusBadRequest)
				logger.Error(rows.Err().Error())
				return
			}

			go st.SendMessage(spam, messageID, conf, logger, ml)
			i++
		}
	}
	if i == 0 {
		logger.Error("There is no Clients with such filter!")
		return
	}
	logger.Info("Send all requests to extern server!")
	return
}

func (st *Postgres) GetPhoneNum(ctx context.Context, clientID int, phone *int) error {
	var phoneStr string
	q := `SELECT phone_num FROM клиент WHERE id = $1;`
	err := st.pool.QueryRow(ctx, q, clientID).Scan(&phoneStr)
	if err != nil {
		return err
	}
	*phone, err = strconv.Atoi(phoneStr)
	if err != nil {
		return err
	}
	return nil
}

func (st *Postgres) SendMessage(spam *Body, messageID int, conf *config.Config, logger *logrus.Logger, ml *validation.MailingList) {
	//deadline for requesting extern API
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(conf.SendMessageTimeout)*time.Millisecond))
	defer cancel()

	url := "https://probe.fbrq.cloud/v1/send/"
	var bearer = "Bearer " + conf.Token

	jsonBody, err := json.Marshal(spam)
	if err != nil {
		logger.Error(err.Error())
	}

	req, err := http.NewRequest("POST", url+strconv.Itoa(messageID), bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error(err.Error())
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	client := &http.Client{}
	for validation.ValidateExternAPIResponse(resp) {
		select {
		case <-ctx.Done():
			logger.Error("[goroutine] MailingList:", ml.ID, " TIMEOUT sendind message to Client:", spam.Id)
			return
		default:
			resp, err = client.Do(req)
			if err != nil {
				continue
			}
			time.Sleep(time.Duration(conf.RequestFreq) * time.Millisecond)
			continue
		}
	}

	q := `UPDATE сообщение
		SET status = true
		WHERE id = $1`
	_, err = st.pool.Exec(context.Background(), q, messageID)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("[goroutine] MailingList:", ml.ID, " successfuly send message to Client:", spam.Id)
	return
}

func (st *Postgres) GeneralStats(ctx context.Context, w http.ResponseWriter) error {
	resp := struct {
		MailingListNumber                 int
		StatusNotFinishedNumberOfMessages int
		StatusFinishedNumberOfMessages    int
	}{}

	q := `SELECT COUNT(*) AS Number
		FROM рассылка`
	err := st.pool.QueryRow(ctx, q).Scan(&resp.MailingListNumber)
	if err != nil {
		return err
	}

	q = `SELECT COUNT(status)
		FROM сообщение
		GROUP BY status`
	rows, err := st.pool.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		if i == 0 {
			rows.Scan(&resp.StatusNotFinishedNumberOfMessages)
		} else {
			rows.Scan(&resp.StatusFinishedNumberOfMessages)
		}
		if rows.Err() != nil {
			return rows.Err()
		}
		i++
	}

	json.NewEncoder(w).Encode(resp)
	return nil
}

func (st *Postgres) GetDetailedStats(ctx context.Context, w http.ResponseWriter, id string) error {
	type Message struct {
		MessageID  int
		CreateDate time.Time
		Status     bool
		ClientID   int
	}
	resp := struct {
		MailingListID      int
		TotalMessagesCount int

		Message    string
		Filter     string
		LaunchDate time.Time
		FinishDate time.Time
		Messages   []Message
	}{}

	strID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("not valid ID format")
	}
	resp.MailingListID = strID

	q := `SELECT exists (SELECT id FROM рассылка WHERE id = $1)`
	var exists bool
	err = st.pool.QueryRow(ctx, q, strID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No such MailListID in database!")
	}

	q = `SELECT message, filter, launch_date, finish_date, сообщение.id, create_date, status, client_id
		FROM рассылка INNER JOIN сообщение
    	on сообщение.mailinglist_id = рассылка.id
		WHERE рассылка.id = $1;`
	rows, err := st.pool.Query(ctx, q, strID)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println(rows.CommandTag().String())

	i := 0
	for rows.Next() {
		resp.Messages = append(resp.Messages, Message{})
		rows.Scan(&resp.Message,
			&resp.Filter,
			&resp.LaunchDate,
			&resp.FinishDate,
			&resp.Messages[i].MessageID,
			&resp.Messages[i].CreateDate,
			&resp.Messages[i].Status,
			&resp.Messages[i].ClientID)
		if rows.Err() != nil {
			return rows.Err()
		}
		i++
	}

	resp.TotalMessagesCount = len(resp.Messages)
	json.NewEncoder(w).Encode(resp)
	return nil
}

func NewPostgres(config *config.Config) (*Postgres, error) {
	postgres := new(Postgres)
	err, pool := NewClient(context.Background(), config)
	if err != nil {
		return nil, err
	}
	postgres.pool = pool
	return postgres, nil
}

func NewClient(ctx context.Context, config *config.Config) (error, *pgxpool.Pool) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Storage.Username,
		config.Storage.Password,
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.Database)
	var pool *pgxpool.Pool

	err := utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, config.Storage.Attempts, 5*time.Second)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return nil, pool
}
