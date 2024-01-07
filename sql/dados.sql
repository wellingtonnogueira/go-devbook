INSERT INTO devbook.usuarios
(nome, nick, email, senha, criadoEm)
VALUES
('Usuario I', 'usuario1', 'usuario1@usuario.com', '$2a$10$5Mo750zv4H4pE0XDTTdsi.vfG2gJ9UcbY3C7yHlk6QA7cxAm0swA2', CURRENT_TIMESTAMP),
('Usuario II', 'usuario2', 'usuario2@usuario.com', '$2a$10$5Mo750zv4H4pE0XDTTdsi.vfG2gJ9UcbY3C7yHlk6QA7cxAm0swA2', CURRENT_TIMESTAMP),
('Usuario III', 'usuario3', 'usuario3@usuario.com', '$2a$10$5Mo750zv4H4pE0XDTTdsi.vfG2gJ9UcbY3C7yHlk6QA7cxAm0swA2', CURRENT_TIMESTAMP);

INSERT INTO devbook.seguidores (usuario_id, seguidor_id)
VALUES
(1, 2),
(3, 1),
(1, 3);

INSERT INTO devbook.publicacoes(titulo, conteudo, autor_id)
VALUES
("Publicação do usuário 1", "Essa é a publicação do usuário 1. Oba!", 1),
("Publicação do usuário 2", "Essa é a publicação do usuário 2. Oba!", 2),
("Publicação do usuário 3", "Essa é a publicação do usuário 3. Oba!", 3);