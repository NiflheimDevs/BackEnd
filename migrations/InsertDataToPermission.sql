ALTER SEQUENCE permission_id_seq RESTART WITH 0;

INSERT INTO "permission" ("name" , "description") VALUES
("ADD_MEMBER", "adds member"),
("REMOVE_MEMEBER", "removes a member"),
("EDIT_INFO", "edit title, bio and ..."),
("BIDDER", "the one who bids"),
("EDIT_NICKNAME", "for teams, it works for positions. for groups and etc for nickname"),
("EDIT_ROLE", "able to change the roles");