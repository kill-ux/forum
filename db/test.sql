
BEGIN TRANSACTION;

CREATE TABLE `categories` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL UNIQUE
);
INSERT INTO categories VALUES(1,'art');
INSERT INTO categories VALUES(2,'programming');
INSERT INTO categories VALUES(3,'news');
INSERT INTO categories VALUES(4,'studying');
INSERT INTO categories VALUES(5,'business');
INSERT INTO categories VALUES(6,'Discussions');
INSERT INTO categories VALUES(7,'Questions');
INSERT INTO categories VALUES(8,'Ideas');
INSERT INTO categories VALUES(9,'Articles');
INSERT INTO categories VALUES(10,'Events');
INSERT INTO categories VALUES(11,'Issues');
INSERT INTO categories VALUES(12,'Others');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('likes',526);
INSERT INTO sqlite_sequence VALUES('comments',75);
INSERT INTO sqlite_sequence VALUES('posts',54);
INSERT INTO sqlite_sequence VALUES('users',18);
INSERT INTO sqlite_sequence VALUES('categories',12);
CREATE TABLE `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` TEXT NOT NULL UNIQUE,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL,
    `image` TEXT NOT NULL,
    `role` TEXT NOT NULL,
    `token` TEXT UNIQUE,
    `token_exp` INTEGER,
    `created_at` DATE NOT NULL
);
INSERT INTO users VALUES(1,'mob','mob@gmail.com','$2a$14$dCdB4Ku4ouoY6eZ4VJSuHueHVJKjDGzqCr0P1ayg5FfnSIsp4e0Ym','res/image/profile.png','user','d3279d80-b2f4-11ef-8cc2-047c1698dc71',1733481405,'2024-11-25 12:32:07.273287208+01:00');
INSERT INTO users VALUES(2,'adnane','adnane@gmail.com','$2a$14$0vdKWBtMAuz.HEtMveaFkudZdKIUUD81CGSyodphDMLwlh76XFjz6','res/image/profile.png','user','afccbf0d-b2f7-11ef-8478-047c1698dc71',1733482634,'2024-11-25 12:33:25.702254098+01:00');
INSERT INTO users VALUES(3,'kabani','kabani@gmail.com','$2a$14$mnTEiJPZBSj03wZPOCeuvOMrng/IHNEBsg6FVXUNnd6WAiJdp6.g2','res/image/profile.png','user','6563f8f9-abf1-11ef-aca6-047c1698dd2b',1732710274,'2024-11-26 12:43:17.240507403+01:00');
INSERT INTO users VALUES(4,'elmir','elmir@gmail.com','$2a$14$nwT0MBrdHXZNGOO5wwClW.5Vf1aLzmi84AdZ7YhSn4SOJUGpIxN/a','res/image/profile.png','user','0acdc7fe-ac00-11ef-906a-047c1698dd2b',1732716565,'2024-11-26 15:09:14.906395375+01:00');
INSERT INTO users VALUES(5,'test','test@gmail.com','$2a$14$hyuu4F4YDcWWJIliK54aQuSkqwjHaRMLhFeD1d3FVmTL5iWBqYVEu','res/image/profile.png','user','9c44f109-ac9e-11ef-8fce-047c1698dd2b',1732784669,'2024-11-26 15:09:31.698351387+01:00');
INSERT INTO users VALUES(6,'mostafa','mm@gmail.com','$2a$14$TKfsUMlVAied.kKCufkjQutp8TU/A46OBU5dgjLZzWmMIkrVa84ji','res/image/profile.png','user','70ef2458-ad7c-11ef-9420-047c1698dd2b',1732879945,'2024-11-28 12:05:10.742846226+01:00');
INSERT INTO users VALUES(7,'Kill','kill@gmail.com','$2a$14$Kxopnseh6Px9LyQMCUxoDe1SDNF3bQ5gz54G8m6ohG1TtcUL4z7Nm','res/image/profile.png','user','a06e149c-b21c-11ef-8e4b-047c1698dc71',1733388549,'2024-11-30 16:15:23.686237357+01:00');
INSERT INTO users VALUES(8,'9696','mob66@gmail.com','$2a$14$D/GyrEjEQ0lf4acNytNw/u49BMf1WbfK4VTLdbSs865i02ZPF3xHy','res/image/profile.png','user',NULL,0,'2024-12-03 18:19:46.365230071+01:00');
INSERT INTO users VALUES(10,'hajji','hajju@gmail.com','$2a$14$RZq/3ZwTtxgJXDuP.Q0Js.QP1MtYnpeg/GGxua.5yyoBY0yY7n3MS','res/image/profile.png','user','712c07ac-b21e-11ef-815a-047c1698dc71',1733389328,'2024-12-04 10:02:00.45444802+01:00');
INSERT INTO users VALUES(11,'error','error@gmail.com','$2a$14$z/KJEyjmRHv19Iaisl.EJesSBQ1t8s0.Ou5tiojdyjP36fE2PBj92','res/image/profile.png','user','2f40d401-b222-11ef-815a-047c1698dc71',1733390936,'2024-12-04 10:04:18.374302261+01:00');
INSERT INTO users VALUES(12,'sqdqsd','sdqds@gg','$2a$14$auyjZEFpQwHaubkS77dSiuZRy/aZllzwUjgZNuVgC1zp/970uTRb.','res/image/profile.png','user',NULL,0,'2024-12-05 10:11:01.053912359+01:00');
INSERT INTO users VALUES(13,'ss','ss@gmail.com','$2a$14$aCRcH5UzQ.aU13FFBCjy/eVm8AIsinfR8j43MN3uZqXpmQFAkexsG','res/image/profile.png','user',NULL,0,'2024-12-05 10:11:25.639460679+01:00');
INSERT INTO users VALUES(14,'','','$2a$14$GmhS5iePNmLuACWhWgPBbu7CC.vOcH7uI5LrWiLk8g/g8ehprvgKK','res/image/profile.png','user','d5ec27bf-b2f4-11ef-8cc2-047c1698dc71',1733481410,'2024-12-05 10:17:20.112939246+01:00');
INSERT INTO users VALUES(15,'ww','ww@gmail.com','$2a$14$WO9ovLjDkEZ.tFQMEQYeje7LHuMmCZ/ru3ATPIogbUHgoia4W6Eay','res/image/profile.png','user','e9b4941e-b2e9-11ef-b386-047c1698dc71',1733476718,'2024-12-05 10:18:10.795068524+01:00');
INSERT INTO users VALUES(16,'sqdsqd','.@j.u.com','$2a$14$pMBdf6LmUhtKLuLng.ysVOD/GikccPUWngIAuEbgiE.bABXKTFFX2','res/image/profile.png','user',NULL,0,'2024-12-05 10:30:09.462733343+01:00');
INSERT INTO users VALUES(17,'kill-ux','kill-ux@gmail.com','$2a$14$Rerhuoo6HWSexqXRI78RJu3rUSfpok0EiLTyOuEhcKKKmwMtJRbqq','res/image/profile.png','user','950f2a3a-b2f7-11ef-8478-047c1698dc71',1733482589,'2024-12-05 11:56:19.615403192+01:00');
INSERT INTO users VALUES(18,'Sphinex','Sphinex@gmail.com','$2a$14$HxSnWNR.a9mpWQuIZFLJTu7CIIVKP1Zu1kU8gBzRK7MtSF9xRBrRa','res/image/profile.png','user','d617ae3a-b2f7-11ef-8478-047c1698dc71',1733482699,'2024-12-05 11:57:08.631499188+01:00');

CREATE TABLE `posts` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `title` TEXT NOT NULL,
    `body` TEXT NOT NULL,
    `image` TEXT,
    `created_at` INTEGER NOT NULL,
    `modified_at` INTEGER NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
INSERT INTO posts VALUES(14,5,'The new languge','The new languge is Fast','',1732640162,1732640162);
INSERT INTO posts VALUES(15,5,'xgxdfg','dgsdg','cc08f8c7-09a7-4c93-9be4-b2f861449854.jpeg',1732641294,1732641294);
INSERT INTO posts VALUES(16,1,'new Post','balbalbal balbalbal balbalbal balbalbal','8162f047-458c-4c73-8bbe-c076d482b92b.jpg',1732726333,1732726333);
INSERT INTO posts VALUES(32,7,'dfgdfgdfgfd',' dfgdgdfg','',1732980892,1732980892);
INSERT INTO posts VALUES(34,1,'Manga of mine',' There is a charachtre that have a white hear with hard parents','b37880d7-7c84-48a0-ad42-8a51b91d749c.jpg',1733219133,1733219133);
INSERT INTO posts VALUES(38,1,'+63',' 6363',NULL,1733246338,1733246338);
INSERT INTO posts VALUES(39,10,'hello','gfgsffgddghfffyh jfgkuysdgf','b6455090-7b08-4be7-84cc-78ab44e20a27.jpg',1733302982,1733302982);
INSERT INTO posts VALUES(40,10,'h',' how are you ','26e0c982-8690-4d3e-96ec-4ce67fe5c1f7.png',1733303117,1733303117);
INSERT INTO posts VALUES(41,11,'404','error  ','',1733303969,1733303969);
INSERT INTO posts VALUES(43,11,'ho is me',' ','',1733304179,1733304179);
INSERT INTO posts VALUES(50,11,'dsssssssssssssss','fdgfdfg','',1733319584,1733319584);
INSERT INTO posts VALUES(51,11,'sdfsdffffffffffffffffffff','fffffffffffffffffffffffffffffffff','',1733319685,1733319685);
INSERT INTO posts VALUES(52,11,'qqqqqqqqqqqqqqqqq','qqqqqqqqqqqqqqqqqqqqq','',1733319730,1733319730);
INSERT INTO posts VALUES(53,17,'Help us do more',replace(replace('"We''ll get right to the point: we''re asking you to help support Khan Academy. We''re a nonprofit that relies on support from people like you. If everyone reading this gives $12 monthly, Khan Academy can continue to thrive for years. Please help keep Khan Academy free, for anyone, anywhere forever.\r\n\r\nGifts from now through Dec 31 will be matched. Give Now!"','\r',char(13)),'\n',char(10)),'7db3243d-d593-48e9-b69e-13c2a2a7069e.jpeg',1733396484,1733396484);
INSERT INTO posts VALUES(54,17,'The most popular programming languages in 2024 (and what that even means) We aggregated data from nine different rankings to produce',replace(replace('We aggregated data from nine different rankings to produce the ZDNET Index of Programming Language Popularity. Here''s which languages came out on top and what to make of this information.\r\n\r\n\r\nWe recently ran a piece that summarized an IEEE study of programming language popularity based on job listings. It definitely fostered some conversation, including some debate about whether the languages IEEE used in its survey were even languages.\r\n\r\nMost of us are familiar with polls and poll results, especially during campaign seasons. Unfortunately, polls have long been proven to be far from accurate. Some polls have a natural bias to one party or the other (not for nefarious reasons, but just based on how they gather their data). Other polls have demographic or psychographic bias. The bottom line is simple: just because the numbers go up in one poll, that doesn''t mean your candidate will win.','\r',char(13)),'\n',char(10)),'508d3bcc-bd3b-4886-b23d-d4d8c333844d.jpeg',1733396564,1733396564);


CREATE TABLE `comments` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER NOT NULL,
    `body` TEXT NOT NULL,
    `created_at` INTEGER NOT NULL,
    `modified_at` INTEGER NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
);
INSERT INTO comments VALUES(45,1,15,'wxa fiha',1732725896,1732725896);
INSERT INTO comments VALUES(46,1,15,'qsdsqd',1732725899,1732725899);
INSERT INTO comments VALUES(47,1,15,'qsdsqd',1732725902,1732725902);
INSERT INTO comments VALUES(48,1,16,'nice subject',1732726352,1732726352);
INSERT INTO comments VALUES(49,1,14,'dsdfsd',1732786242,1732786242);
INSERT INTO comments VALUES(50,1,6,'wx bik',1733240165,1733240165);
INSERT INTO comments VALUES(51,1,34,'wx bik',1733240809,1733240809);
INSERT INTO comments VALUES(52,2,34,'don''t read it',1733242784,1733242784);
INSERT INTO comments VALUES(53,1,34,'qsdqsd',1733243396,1733243396);
INSERT INTO comments VALUES(56,11,40,'fuck you',1733303168,1733303168);
INSERT INTO comments VALUES(57,11,40,'hjkhjk',1733303580,1733303580);
INSERT INTO comments VALUES(58,11,40,'yuiuyi',1733303583,1733303583);
INSERT INTO comments VALUES(59,11,40,'yuiyuiyu',1733303586,1733303586);
INSERT INTO comments VALUES(60,11,40,'yuiyui',1733303588,1733303588);
INSERT INTO comments VALUES(61,11,40,'yuikyuiyuh',1733303590,1733303590);
INSERT INTO comments VALUES(62,11,40,'yiuyi',1733303625,1733303625);
INSERT INTO comments VALUES(63,2,43,'wx bik',1733319319,1733319319);
INSERT INTO comments VALUES(64,2,43,'dswfqsf',1733319323,1733319323);
INSERT INTO comments VALUES(65,2,43,'OnloadBody',1733319329,1733319329);
INSERT INTO comments VALUES(66,2,43,'OnloadBody 4',1733319334,1733319334);
INSERT INTO comments VALUES(67,2,43,'OnloadBody 2',1733319343,1733319343);
INSERT INTO comments VALUES(68,2,43,'OnloadBody 4',1733319346,1733319346);
INSERT INTO comments VALUES(69,2,43,'404',1733319354,1733319354);
INSERT INTO comments VALUES(70,11,43,'<a href=''/ss''>we</a>',1733319384,1733319384);
INSERT INTO comments VALUES(71,11,43,'dsdfdsf',1733319451,1733319451);
INSERT INTO comments VALUES(72,11,43,'sdfdsf',1733319453,1733319453);
INSERT INTO comments VALUES(74,17,53,'kayn',1733396580,1733396580);
INSERT INTO comments VALUES(75,18,53,'not good',1733396588,1733396588);
CREATE TABLE `posts_categories` (
    `post_id` INTEGER NOT NULL,
    `category_id` INTEGER NOT NULL,
    PRIMARY KEY (`category_id`, `post_id`),
    FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE
);
INSERT INTO posts_categories VALUES(14,2);
INSERT INTO posts_categories VALUES(14,3);
INSERT INTO posts_categories VALUES(14,6);
INSERT INTO posts_categories VALUES(15,4);
INSERT INTO posts_categories VALUES(15,7);
INSERT INTO posts_categories VALUES(15,10);
INSERT INTO posts_categories VALUES(16,2);
INSERT INTO posts_categories VALUES(16,6);
INSERT INTO posts_categories VALUES(16,9);
INSERT INTO posts_categories VALUES(17,9);
INSERT INTO posts_categories VALUES(32,2);
INSERT INTO posts_categories VALUES(32,3);
INSERT INTO posts_categories VALUES(32,4);
INSERT INTO posts_categories VALUES(32,5);
INSERT INTO posts_categories VALUES(32,6);
INSERT INTO posts_categories VALUES(32,7);
INSERT INTO posts_categories VALUES(32,8);
INSERT INTO posts_categories VALUES(32,9);
INSERT INTO posts_categories VALUES(32,10);
INSERT INTO posts_categories VALUES(32,11);
INSERT INTO posts_categories VALUES(32,12);
INSERT INTO posts_categories VALUES(34,2);
INSERT INTO posts_categories VALUES(34,3);
INSERT INTO posts_categories VALUES(34,4);
INSERT INTO posts_categories VALUES(34,8);
INSERT INTO posts_categories VALUES(38,1);
INSERT INTO posts_categories VALUES(38,2);
INSERT INTO posts_categories VALUES(38,3);
INSERT INTO posts_categories VALUES(38,5);
INSERT INTO posts_categories VALUES(40,1);
INSERT INTO posts_categories VALUES(40,2);
INSERT INTO posts_categories VALUES(40,3);
INSERT INTO posts_categories VALUES(40,6);
INSERT INTO posts_categories VALUES(40,7);
INSERT INTO posts_categories VALUES(40,8);
INSERT INTO posts_categories VALUES(40,10);
INSERT INTO posts_categories VALUES(40,11);
INSERT INTO posts_categories VALUES(43,7);
INSERT INTO posts_categories VALUES(50,5);
INSERT INTO posts_categories VALUES(51,1);
INSERT INTO posts_categories VALUES(52,11);
INSERT INTO posts_categories VALUES(53,2);
INSERT INTO posts_categories VALUES(53,12);
INSERT INTO posts_categories VALUES(54,12);
CREATE TABLE `likes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` INTEGER NOT NULL,
    `post_id` INTEGER ,
    `comment_id` INTEGER ,
    `like` BOOLEAN NOT NULL ,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE
);
INSERT INTO likes VALUES(120,5,15,NULL,1);
INSERT INTO likes VALUES(201,4,8,NULL,1);
INSERT INTO likes VALUES(205,4,15,NULL,1);
INSERT INTO likes VALUES(247,1,NULL,45,0);
INSERT INTO likes VALUES(251,1,NULL,46,0);
INSERT INTO likes VALUES(252,1,NULL,47,0);
INSERT INTO likes VALUES(254,1,NULL,48,0);
INSERT INTO likes VALUES(261,1,3,NULL,0);
INSERT INTO likes VALUES(263,1,5,NULL,1);
INSERT INTO likes VALUES(266,1,15,NULL,1);
INSERT INTO likes VALUES(281,1,32,NULL,1);
INSERT INTO likes VALUES(282,1,6,NULL,1);
INSERT INTO likes VALUES(284,1,14,NULL,1);
INSERT INTO likes VALUES(285,1,16,NULL,1);
INSERT INTO likes VALUES(286,1,34,NULL,1);
INSERT INTO likes VALUES(287,1,NULL,51,1);
INSERT INTO likes VALUES(288,10,32,NULL,0);
INSERT INTO likes VALUES(289,10,16,NULL,1);
INSERT INTO likes VALUES(290,11,39,NULL,0);
INSERT INTO likes VALUES(291,11,40,NULL,1);
INSERT INTO likes VALUES(293,10,NULL,56,0);
INSERT INTO likes VALUES(305,2,NULL,57,0);
INSERT INTO likes VALUES(307,2,NULL,56,1);
INSERT INTO likes VALUES(310,2,NULL,58,0);
INSERT INTO likes VALUES(313,2,NULL,69,1);
INSERT INTO likes VALUES(316,11,NULL,65,1);
INSERT INTO likes VALUES(317,2,NULL,67,1);
INSERT INTO likes VALUES(318,11,NULL,68,0);
INSERT INTO likes VALUES(319,11,NULL,70,1);
INSERT INTO likes VALUES(320,2,NULL,68,1);
INSERT INTO likes VALUES(323,2,NULL,66,1);
INSERT INTO likes VALUES(435,2,43,NULL,0);
INSERT INTO likes VALUES(481,2,38,NULL,1);
INSERT INTO likes VALUES(482,2,14,NULL,1);
INSERT INTO likes VALUES(483,2,15,NULL,1);
INSERT INTO likes VALUES(484,2,16,NULL,1);
INSERT INTO likes VALUES(512,11,52,NULL,0);
INSERT INTO likes VALUES(517,2,40,NULL,1);
INSERT INTO likes VALUES(518,2,34,NULL,0);
INSERT INTO likes VALUES(520,2,50,NULL,1);
INSERT INTO likes VALUES(521,2,39,NULL,1);
INSERT INTO likes VALUES(522,2,41,NULL,1);
INSERT INTO likes VALUES(523,18,53,NULL,0);
INSERT INTO likes VALUES(524,17,53,NULL,1);
INSERT INTO likes VALUES(525,17,NULL,73,1);
INSERT INTO likes VALUES(526,17,NULL,74,1);

COMMIT;
