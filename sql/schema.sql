CREATE TABLE foo (
  id INT PRIMARY KEY,
  name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS nail_users (
    id SERIAL NOT NULL,
    username VARCHAR ( 255 ) UNIQUE NOT NULL,
    password VARCHAR ( 255 ) NOT NULL,
    created_on TIMESTAMP,
    last_login TIMESTAMP
)

CREATE TABLE IF NOT EXISTS nail_colors (
    sku INTEGER,
    name VARCHAR ( 255 ),
    emoji VARCHAR ( 255 ),
    cssname VARCHAR ( 255 )
)
-- they don't match
-- INSERT INTO nail_colors (sku, name, emoji, cssname) VALUES
-- (0, 'black', '⬛', 'black'),
-- (1, 'red', '🟥', 'red'),
-- (2, 'turquoise', '🟦', 'turquoise'),
-- (3, 'yellow', '🟨', 'yellow'),
-- (4, 'green', '🟩', 'green'),
-- (5, 'purple', '🟪', 'purple'),
-- (6, 'white', '⬜', 'white'),
-- (7, 'grey', '🔲', 'grey'),
-- (8, 'gold', '🥇', 'gold'),
-- (9, 'rose', '🌹', 'pink')

-- insert colors = [
-- ['black', '⬛', 'black'],
-- ['red', '🟥', 'red'],
-- ['TQ', '🟦', 'turquoise'],
-- ['yellow', '🟨', 'yellow'],
-- ['green', '🟩', 'green'],
-- ['purple', '🟪', 'purple'],
-- ['white', '⬜', 'white'],
-- ['grey', '🔲', 'grey'],
-- ['gold', '🥇', 'gold'],
-- ['rose', '🌹', 'pink']
-- ]

CREATE TABLE IF NOT EXISTS nail_parts (
    name VARCHAR ( 255 ) NOT NULL,
    size VARCHAR ( 255 ) NOT NULL,
    color VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_items (
    name VARCHAR ( 255 ) NOT NULL,
    size VARCHAR ( 255 ) NOT NULL,
    a_color VARCHAR ( 255 ),
    b_color VARCHAR ( 255 ),
    c_color VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_types (
    name VARCHAR ( 255 ),
    sku INTEGER
)
-- insert types = [
-- ['Laser Cut', '0'],
-- ['Tee Shirt', '1'],
-- ['Tank Top', '2'],
-- ['Hoodie', '3'],
-- ['Screen Print', '10'],
-- ['Greeting Card', '11']
-- ]

CREATE TABLE IF NOT EXISTS nail_shirts (
    nombre VARCHAR ( 255 ) NOT NULL,
    a VARCHAR ( 255 ),
    b VARCHAR ( 255 ),
    c VARCHAR ( 255 ),
    backs VARCHAR ( 255 ),
    sku INTEGER
)

CREATE TABLE IF NOT EXISTS nail_boxes (
    name VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_boxprod (
    name VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_boxused (
    name VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_projections (
    name VARCHAR ( 255 ) NOT NULL,
    size VARCHAR ( 255 ) NOT NULL,
    a_color VARCHAR ( 255 ),
    b_color VARCHAR ( 255 ),
    c_color VARCHAR ( 255 ),
    qty INTEGER,
    cycle INTEGER,
    sku BIGINT
)

CREATE TABLE IF NOT EXISTS nail_queueItems (
    name VARCHAR ( 255 ) NOT NULL,
    size VARCHAR ( 255 ) NOT NULL,
    a_color VARCHAR ( 255 ),
    b_color VARCHAR ( 255 ),
    c_color VARCHAR ( 255 ),
    qty INTEGER,
    cycle INTEGER,
    sku BIGINT
)

CREATE TABLE IF NOT EXISTS nail_queueParts (
    name VARCHAR ( 255 ) NOT NULL,
    size VARCHAR ( 255 ) NOT NULL,
    color VARCHAR ( 255 ),
    qty INTEGER
)

CREATE TABLE IF NOT EXISTS nail_cycles (
    id SERIAL UNIQUE NOT NULL,
    name VARCHAR (255),
    created_on TIMESTAMP,
    current BOOL
)
