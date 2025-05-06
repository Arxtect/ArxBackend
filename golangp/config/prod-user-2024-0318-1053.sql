INSERT INTO public.users
(id, name, email, password, role, provider, photo, verification_code, password_reset_token, password_reset_at, verified, balance, created_at, updated_at)
VALUES
    ('85ae43c0-b9f6-464f-b54a-29701d7ede5d', 'awOThJStbo', 'UnZUCkGgBc@example.com', '$2a$10$Gs34LJptLJrnHhmtbnIEzO96nCDBsepMfc82eMEjS5aslTnXYi3c2', 'user', 'local', NULL, NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-05 15:56:58.939814+00', '2024-02-05 15:56:58.939815+00'),
    ('27cd3aa3-0d4a-43ce-8fb1-3508459cbda0', 'WYhiTmfOpE', 'EPTYnVJcLS@example.com', '$2a$10$Gs34LJptLJrnHhmtbnIEzO96nCDBsepMfc82eMEjS5aslTnXYi3c2', 'user', 'local', NULL, NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-05 15:56:58.941793+00', '2024-02-05 15:56:58.941793+00'),
    ('93d08112-9b05-44ca-bc47-1ad098e6efb9', 'dQcIFzWbHH', 'lUzWPLZKcn@example.com', '$2a$10$Gs34LJptLJrnHhmtbnIEzO96nCDBsepMfc82eMEjS5aslTnXYi3c2', 'user', 'local', NULL, NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-05 15:56:58.942781+00', '2024-02-05 15:56:58.942781+00'),
    ('192bd05f-6625-4c78-af09-18d454d736a8', 'llSAHFIZXb', 'qjyeYJNCJn@example.com', '$2a$10$Gs34LJptLJrnHhmtbnIEzO96nCDBsepMfc82eMEjS5aslTnXYi3c2', 'user', 'local', NULL, NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-05 15:56:58.943709+00', '2024-02-05 15:56:58.94371+00'),
    ('6382203f-bb5f-4afc-822b-e72c6047cbf5', 'UeUdaAtAAc', 'ThANEphZXm@example.com', '$2a$10$Gs34LJptLJrnHhmtbnIEzO96nCDBsepMfc82eMEjS5aslTnXYi3c2', 'user', 'local', NULL, NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-05 15:56:58.944623+00', '2024-02-05 15:56:58.944623+00'),
    ('49e0b32e-c516-473c-a26f-1744751e3c35', 'devin', 'zhengdevin10@gmail.com', '$2a$10$xX5aWvW35I4TF10LklGqK.AsiisP11JBvsgBOhTQcampjzadzXKJa', 'user', 'local', 'test', 'Vk4yT29PaUU4TDVJNEV4ajlKUVo=', NULL, '0001-01-01 00:00:00+00', 'f', 0, '2024-02-06 02:21:37.397286+00', '2024-02-06 02:21:37.402185+00'),
    ('862c8105-f092-42f5-b484-cb689e7a2f9e', 'devin.', '2473023641@qq.com', '$2a$10$aBwvLjzEmhGgTI6XWYAEFe/KgGn5XVsWD9pUiz5QbKPNEkEhHko2m', 'user', 'local', 'test', NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-06 02:23:20.17292+00', '2024-02-06 02:24:03.661482+00'),
    ('00bd22c5-746d-4e97-855f-e72020d46a3a', 'fucewei', '1457310354@qq.com', '$2a$10$2S3lbbwYI86W8nA.Dkw5v.hRCYo9/5BZ5JKJsvHXWNQuBlW.Wb4fC', 'user', 'local', 'test', NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-06 02:38:44.454336+00', '2024-02-06 02:40:28.419829+00'),
    ('814f9673-0103-4440-b975-fb0913855507', 'Admin Admin', 'tzion@open2any.tools', '$2a$10$Quh8rFc1ilZLA4pLe3rnQe1jMDbUnRVC1eFH71XI3xQH2as41Q6YG', 'admin', 'local', 'test', NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-29 09:00:35.848907+00', '2024-02-29 09:00:35.848907+00'),
    -- ('1f2a0352-686d-4a1b-8a90-e1054472cccd', 'Admin Admin', 'admin@open2any.tools', '$2a$10$Quh8rFc1ilZLA4pLe3rnQe1jMDbUnRVC1eFH71XI3xQH2as41Q6YG', 'admin', 'local', 'test', NULL, NULL, '0001-01-01 00:00:00+00', 't', 0, '2024-02-29 09:00:35.853023+00', '2024-02-29 09:00:35.853023+00'),
    ('00b5af28-70a5-449e-93b9-fb8143f04e36', 'ethanshanyu', 'ethanshanyu@gmail.com', '$2a$10$kDmN.4D5f9KO2HG7SjgS8unEVLfkikkzOgXISxhe.r/.PVsyuB5Py', 'user', 'local', 'test', 'ekZJemZBU051VkZmaVJYcDdtNmQ=', NULL, '0001-01-01 00:00:00+00', 'f', 0, '2024-03-03 10:57:00.429582+00', '2024-03-03 10:57:00.433868+00');



INSERT INTO public.tags
(id,created_at,updated_at,name)
VALUES
    ('6db5d007-9202-4f18-ab17-3d4cbd79e99f','2024-01-25 15:16:25.158+08',NULL,'astronomy'),
    ('d090d38e-9219-4d13-b543-9605e73381cf','2024-01-19 05:54:44.041+08',NULL,'universe'),
    ('b710a0f0-b428-438c-bf6b-2486322378d2','2024-02-23 01:07:10.217+08',NULL,'machinery'),
    ('b615888b-0003-40a6-a345-6209ae61889b','2024-02-13 18:14:26.537+08',NULL,'internet'),
    ('676ebce1-dd6a-4929-b0a1-4d23408ca0e8','2024-01-15 05:42:48.633+08',NULL,'physics'),
    ('e3829439-3d21-4e4a-b4f7-cb07b60783f0','2024-02-24 04:31:59.193+08',NULL,'romantic'),
    ('0e47ca1e-07c3-4680-9508-ae7ee55e0a9a','2024-02-17 14:44:46.537+08',NULL,'scientific'),
    ('b7f55dec-05cc-4cb0-ab70-8d66e6315407','2024-01-14 01:38:50.825+08',NULL,'plant'),
    ('b6f7c402-19e7-411c-a87d-62915326f88d','2024-02-24 06:49:37.865+08',NULL,'organism');