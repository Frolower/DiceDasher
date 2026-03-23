\connect character_db

INSERT INTO tes.car_stats (
    type,
    code,
    passengers,
    maneuverability,
    speed,
    hull,
    armor,
    cost
) VALUES
      ('Horse', 'horse', 1, NULL, 2, 3, NULL, 1000),
      ('Wagon', 'wagon', 4, 1, 2, 6, 2, 200),
      ('Bicycle', 'bicycle', NULL, 1, 2, 3, NULL, 300),
      ('Motorcycle', 'motorcycle', 1, 2, 3, 4, 2, 5000),
      ('Dirt bike', 'dirt_bike', NULL, 3, 2, 3, NULL, 3000),
      ('2WD Car', '2wd_car', 4, 2, 3, 6, 4, 10000),
      ('4WD Car', '4wd_car', 4, 3, 2, 6, 4, 15000),
      ('Pickup Truck', 'pickup_truck', 5, 2, 2, 8, 4, 20000),
      ('Van', 'van', 7, 2, 2, 8, 4, 25000),
      ('Light Truck', 'light_truck', 14, 1, 2, 12, 4, 30000),
      ('Heavy Truck', 'heavy_truck', 16, 1, 2, 14, 4, 50000),
      ('Bus', 'bus', 50, 1, 2, 12, 4, 40000),
      ('Rowboat', 'rowboat', 4, 1, 1, 5, 2, 500),
      ('Small Sailing Boat', 'small_sailing_boat', 7, 1, 2, 6, 2, 5000),
      ('Small Motorboat', 'small_motorboat', 7, 2, 3, 5, 2, 12000),
      ('Helicopter', 'helicopter', 4, 3, 4, 6, 3, 200000),
      ('Light airplane', 'light_airplane', 3, 2, 4, 5, 3, 100000),
      ('Small commercial drone ship', 'small_commercial_drone_ship', 10, 2, 4, 9, 3, 200000),
      ('Military drone ship', 'military_drone_ship', NULL, 3, 5, 12, 8, NULL);