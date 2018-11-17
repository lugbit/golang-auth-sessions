DROP DATABASE IF EXISTS userAuthDB;
CREATE DATABASE userAuthDB;
USE userAuthDB;

/* Registered users table */
CREATE TABLE tblUsers (
	fldID 		 BIGINT AUTO_INCREMENT PRIMARY KEY,
	fldFirstName VARCHAR(50) NOT NULL,
	fldLastName  VARCHAR(50) NOT NULL,
	fldEmail     VARCHAR(50) NOT NULL UNIQUE,
	fldPassword  VARCHAR(255) NOT NULL,
    fldCreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    fldUpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
);

/* Stores tokens used to verify a user's email address. */
CREATE TABLE tblActivationTokens (
	fldToken 		   VARCHAR(255) NOT NULL UNIQUE,
    fldTokenDateIssued TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    fldTokenDateUsed   TIMESTAMP DEFAULT NULL,
    fldIsActivated     BOOL NOT NULL DEFAULT FALSE,
    fldFKUserID        BIGINT NOT NULL,
    PRIMARY KEY(fldToken, fldFKUserID),
    FOREIGN KEY(fldFKUserID) REFERENCES tblUsers(fldID) ON DELETE CASCADE
);

/* Stores user session identifiers */
CREATE TABLE tblSessions (
	fldSessionID VARCHAR(255) NOT NULL UNIQUE,
    fldCreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    fldLastActive TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    fldFKUserID  BIGINT NOT NULL,
    PRIMARY KEY(fldFKUserID, fldSessionID),
    FOREIGN KEY(fldFKUserID) REFERENCES tblUsers(fldID) ON DELETE CASCADE
);
