/* userLogin queries */

SELECT * FROM tblUsers;

/* User insertions */
INSERT INTO tblUsers(fldFirstName, fldLastName, fldEmail, fldPassword) 
			VALUES ("Marck", "Munoz", "Marck527@gmail.com", "Password1"),
				   ("Colin", "Stewart", "Colinps@gmail.com", "hunter2"),
                   ("Brayden", "Gravestock", "braydengravestock@gmail.com", "pineapple");


SELECT * FROM tblActivationTokens;

INSERT INTO tblActivationTokens(fldToken, fldFKUserID) 
			VALUES ("t1-fdksnfdjsndjs223", 1);
                   
                   
-- Selects email address
SELECT 
	fldEmail 
FROM tblUsers 
WHERE fldEmail = ?;

-- Insert user
INSERT INTO tblUsers(fldFirstName, fldLastName, fldEmail, fldPassword) 
				VALUES (?, ?, ?, ?, ?, ?);
                
BEGIN;
INSERT INTO tblUsers(fldFirstName, fldLastName, fldEmail, fldPassword)
  VALUES('Marck', 'Munoz', 'Marck527@gmail.com', 'Password1');
INSERT INTO tblActivationTokens (fldToken,fldFKUserID) 
  VALUES('t12-dsdsdsdvsgds4545', LAST_INSERT_ID());
COMMIT;

-- Select users and tokens
SELECT 
	*
FROM tblUsers
INNER JOIN tblActivationTokens
	ON tblUsers.fldID = tblActivationTokens.fldFKUserID;
    
-- Check if token is expired
SELECT 
	TIMESTAMPDIFF(MINUTE, fldTokenDateIssued, NOW()) AS tokenActiveInMinutes
FROM tblActivationTokens
WHERE fldToken = '659b9635-b5a5-4006-b6a5-017c09c4b852';

-- Mark token as used 
UPDATE tblUsers
INNER JOIN tblActivationTokens
	ON tblActivationTokens.fldFKUserID = tblUsers.fldID
SET fldTokenDateUsed = NOW()
WHERE fldToken = 'a3311ee5-3b6e-4fbb-9441-adb7ba0d5f8a';

-- Check if token is alreadsy used
SELECT 
	fldTokenDateUsed
FROM tblActivationTokens
WHERE fldToken = '51ffcb15-0af3-472a-83e2-70f3efe435ba' AND fldTokenDateUsed IS NOT NULL;

-- Clear token date used
UPDATE tblActivationTokens
SET fldTokenDateUsed = NULL
WHERE fldFKUserID = 1;

-- Activate user
UPDATE tblActivationTokens
SET fldIsActivated = true
WHERE fldToken = 'New token mang';

-- Update token
UPDATE tblUsers
INNER JOIN tblActivationTokens
	ON tblActivationTokens.fldFKUserID = tblUsers.fldID
SET fldToken = "New token mang", fldTokenDateIssued = NOW(), fldTokenDateUsed = NULL
	WHERE fldEmail = "Marck527@gmail.com";
    
-- Returns true if user is activated
SELECT 
	fldIsActivated
FROM tblUsers
INNER JOIN tblActivationTokens
	ON tblActivationTokens.fldFKUserID = tblUsers.fldID
WHERE fldEmail = "Marck527@gmail.com";

-- Unactivate user
UPDATE tblActivationTokens
INNER JOIN tblUsers
	ON tblUsers.fldID = tblActivationTokens.fldFKUserID
SET fldIsActivated = 0
WHERE fldEmail = "Marck527@gmail.com";

-- Check if user has a session
SELECT * FROM tblSessions WHERE fldSessionID = ?

-- Create new session
INSERT INTO tblSessions(fldSessionID, fldFKUserID) VALUES (?, ?)

-- Delete sessions for a user
DELETE FROM tblSessions WHERE fldFKUserID = ?;

-- Search for username
SELECT fldUsername FROM tblUser WHERE fldUsername = ?;

-- Select sessions
SELECT * FROM tblSessions;

-- Get user by email
SELECT 
	fldID,
	fldFirstName,
    fldLastName,
    fldEmail
FROM tblUsers
WHERE fldEmail = "Marck527@gmail.com";

-- Get user by session id
SELECT
	tblUsers.fldID,
    tblUsers.fldFirstName,
    tblUsers.fldLastName,
    tblUsers.fldEmail,
    tblUsers.fldPassword
FROM tblUsers
INNER JOIN tblSessions
	ON tblSessions.fldFKUserID = tblUsers.fldID
WHERE tblSessions.tblSessionID = ?;

-- Returns user id based on sessionid
SELECT
	fldFKUserID
FROM tblSessions
WHERE fldSessionID = ?;

-- Check if user account is activated
SELECT
	fldisActivated
FROM tblUsers
INNER JOIN	tblActivationTokens
	ON tblActivationTokens.fldFKUserID = tblUsers.fldID
WHERE tblUsers.fldEmail = "Marck527@gmail.com";

-- Updates the user's password
UPDATE tblUsers
INNER JOIN tblSessions
	ON tblSessions.fldFKUserID = tblUsers.fldID
SET tblUsers.fldPassword = "Test1"
WHERE tblSessions.fldSessionID = '3d1a3a49-41d2-4914-93b7-8e518c7eaa5d';

-- JOIN ALL!!
SELECT 
	tblUsers.fldID,
    tblUsers.fldFirstName,
    tblUsers.fldLastName,
    tblUsers.fldEmail,
    tblUsers.fldPassword,
    tblActivationTokens.fldToken,
    tblActivationTokens.fldTokenDateIssued,
    tblActivationTokens.fldTokenDateUsed,
    tblActivationTokens.fldIsActivated,
    tblSessions.fldSessionID,
    tblSessions.fldCreatedAt,
    tblSessions.fldLastActive
FROM
	tblUsers
INNER JOIN tblActivationTokens
	ON tblActivationTokens.fldFKUserID = tblUsers.fldID
LEFT JOIN tblSessions
	ON  tblSessions.fldFKUserID = tblUsers.fldID;
    
-- Update session last active
UPDATE tblSessions
SET fldLastActive = CURRENT_TIMESTAMP()
WHERE fldSessionID = '337fd9e1-5d06-470b-a812-61cf80707b04';

-- Delete session
DELETE FROM tblSessions
WHERE fldSessionID = '23b1a18f-f34c-4241-b564-7ba16d72bcdb';

DELETE FROM tblSessions WHERE fldFKUserID = 1;
SELECT * FROM tblSessions;





