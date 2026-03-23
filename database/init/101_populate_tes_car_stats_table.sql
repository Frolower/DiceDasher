CREATE TABLE tes.car_stats (
                               id SERIAL PRIMARY KEY,
                               type VARCHAR(255) NOT NULL,
                               passengers INT NOT NULL,
                               maneuverability INT NOT NULL,
                               speed INT NOT NULL,
                               hull INT NOT NULL,
                               armor INT NOT NULL,
                               cost INT NOT NULL
);

\connect character_db

INSERT INTO tes.car_stats (
    type,
    passengers,
    maneuverability,
    speed,
    hull,
    armor,
    cost
) VALUES
      ('Horse', 1, 0, 2, 3, 0, 1000),
      ('Wagon', 4, 1, 2, 6, 2, 200),
      ('Bicycle', 0, 1, 2, 3, 0, 300),
      ('Motorcycle', 1, 2, 3, 4, 2, 5000),
      ('Dirt bike', 0, 3, 2, 3, 0, 3000),
      ('2WD Car', 4, 2, 3, 6, 4, 10000),
      ('4WD Car', 4, 3, 2, 6, 4, 15000),
      ('Pickup Truck', 5, 2, 2, 8, 4, 20000),
      ('Van', 7, 2, 2, 8, 4, 25000),
      ('Light Truck', 14, 1, 2, 12, 4, 30000),
      ('Heavy Truck', 16, 1, 2, 14, 4, 50000),
      ('Bus', 50, 1, 2, 12, 4, 40000),
      ('Rowboat', 4, 1, 1, 5, 2, 500),
      ('Small Sailing Boat', 7, 1, 2, 6, 2, 5000),
      ('Small Motorboat', 7, 2, 3, 5, 2, 12000),
      ('Helicopter', 4, 3, 4, 6, 3, 200000),
      ('Light airplane', 3, 2, 4, 5, 3, 100000),
      ('Small commercial drone ship', 10, 2, 4, 9, 3, 200000),
      ('Military drone ship', 0, 3, 5, 12, 8, 0);