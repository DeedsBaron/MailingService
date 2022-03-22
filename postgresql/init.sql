CREATE DATABASE spam;
\c spam;
CREATE TABLE рассылка (
                          id SERIAL PRIMARY KEY,
                          launch_date TIMESTAMP NOT NULL,
                          message TEXT UNIQUE NOT NULL,
                          filter TEXT NOT NULL,
                          finish_date TIMESTAMP NOT NULL
);
CREATE TABLE клиент (
                        id SERIAL PRIMARY KEY,
                        phone_num TEXT UNIQUE NOT NULL,
                        mobile_code TEXT NOT NULL,
                        tag TEXT NOT NULL,
                        timezone_abbr TEXT
);
CREATE TABLE сообщение (
                        id SERIAL PRIMARY KEY,
                        create_date TIMESTAMP ,
                        status bool,
                        mailinglist_id INT,
                        client_id INT,
                        CONSTRAINT fk_spam_id
                            FOREIGN KEY(mailinglist_id)
                                REFERENCES рассылка(id)
                                ON DELETE CASCADE,
                        CONSTRAINT fk_client_id
                            FOREIGN KEY(client_id)
                                REFERENCES клиент(id)
                                ON DELETE CASCADE
);
