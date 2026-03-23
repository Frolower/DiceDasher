\connect character_db

-- Weapons

INSERT INTO tes.items (
    code,
    name,
    item_kind,
    data
) VALUES
      (
          'unarmed',
          'Unarmed',
          'weapon',
          '{
            "bonus": null,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": null,
            "tags": []
          }'::jsonb
      ),
      (
          'improvised_club',
          'Improvised club',
          'weapon',
          '{
            "bonus": 1,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": null,
            "tags": []
          }'::jsonb
      ),
      (
          'knife',
          'Knife',
          'weapon',
          '{
            "bonus": 1,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": 25,
            "tags": []
          }'::jsonb
      ),
      (
          'baseball_bat',
          'Baseball bat',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": 50,
            "tags": []
          }'::jsonb
      ),
      (
          'axe',
          'Axe',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": 100,
            "tags": []
          }'::jsonb
      ),
      (
          'chainsaw',
          'Chainsaw',
          'weapon',
          '{
            "bonus": 1,
            "damage_value": 3,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": 250,
            "tags": []
          }'::jsonb
      ),
      (
          'rock',
          'Rock',
          'weapon',
          '{
            "bonus": null,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "medium",
            "price": null,
            "tags": []
          }'::jsonb
      ),
      (
          'taser',
          'Taser',
          'weapon',
          '{
            "bonus": 3,
            "damage_value": null,
            "damage_kind": "effect",
            "range_min": "engaged",
            "range_max": "short",
            "price": 500,
            "tags": ["no_direct_damage"],
            "note": "Target needs to roll for Strength with -2 dice or lose their next turn."
          }'::jsonb
      ),
      (
          'derringer',
          'Derringer',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "short",
            "price": 250,
            "tags": []
          }'::jsonb
      ),
      (
          'handgun',
          'Handgun',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 300,
            "tags": []
          }'::jsonb
      ),
      (
          'magnum_revolver',
          'Magnum revolver',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 3,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 400,
            "tags": []
          }'::jsonb
      ),
      (
          'crossbow',
          'Crossbow',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "medium",
            "range_max": "long",
            "price": 150,
            "tags": []
          }'::jsonb
      ),
      (
          'hunting_rifle',
          'Hunting rifle',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "medium",
            "range_max": "long",
            "price": 500,
            "tags": []
          }'::jsonb
      ),
      (
          'shotgun',
          'Shotgun',
          'weapon',
          '{
            "bonus": 3,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 350,
            "tags": []
          }'::jsonb
      ),
      (
          'sniper_rifle',
          'Sniper rifle',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "medium",
            "range_max": "extreme",
            "price": null,
            "tags": ["not_commercial"]
          }'::jsonb
      ),
      (
          'submachinegun',
          'Submachinegun',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": null,
            "tags": ["not_commercial"]
          }'::jsonb
      ),
      (
          'assault_rifle',
          'Assault rifle',
          'weapon',
          '{
            "bonus": 3,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "long",
            "price": null,
            "tags": ["not_commercial"]
          }'::jsonb
      ),
      (
          'heavy_machinegun',
          'Heavy machinegun',
          'weapon',
          '{
            "bonus": 3,
            "damage_value": 3,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "long",
            "price": null,
            "tags": ["not_commercial"]
          }'::jsonb
      ),
      (
          'neodymium_cannon',
          'Neodymium cannon',
          'weapon',
          '{
            "bonus": null,
            "damage_value": 4,
            "damage_kind": "direct",
            "range_min": "medium",
            "range_max": "long",
            "price": null,
            "tags": ["not_commercial", "neurocaster_only"],
            "note": "Uses Network attribute as a gear bonus."
          }'::jsonb
      ),
      (
          'molotov_cocktail',
          'Molotov cocktail',
          'weapon',
          '{
            "bonus": null,
            "damage_value": 6,
            "damage_kind": "blast",
            "range_min": "engaged",
            "range_max": "medium",
            "price": null,
            "tags": ["no_direct_damage", "explosive"]
          }'::jsonb
      ),
      (
          'hand_grenade',
          'Hand grenade',
          'weapon',
          '{
            "bonus": null,
            "damage_value": 8,
            "damage_kind": "blast",
            "range_min": "engaged",
            "range_max": "medium",
            "price": null,
            "tags": ["no_direct_damage", "not_commercial", "explosive"]
          }'::jsonb
      ),
      (
          'mortar',
          'Mortar',
          'weapon',
          '{
            "bonus": 2,
            "damage_value": 12,
            "damage_kind": "blast",
            "range_min": "medium",
            "range_max": "extreme",
            "price": null,
            "tags": ["no_direct_damage", "not_commercial", "explosive"]
          }'::jsonb
      );

-- Armor

INSERT INTO tes.items (
    code,
    name,
    item_kind,
    data
) VALUES
      (
          'soft_vest',
          'Soft vest',
          'armor',
          '{
            "armor_level": 2,
            "agility_modifier": -1,
            "price": 150
          }'::jsonb
      ),
      (
          'plate_vest',
          'Plate vest',
          'armor',
          '{
            "armor_level": 4,
            "agility_modifier": -2,
            "price": 300
          }'::jsonb
      ),
      (
          'riot_gear',
          'Riot gear',
          'armor',
          '{
            "armor_level": 6,
            "agility_modifier": -3,
            "price": 500
          }'::jsonb
      );

-- Neurocasters

INSERT INTO tes.items (
    code,
    name,
    item_kind,
    data
) VALUES
      (
          'stimulus_tle_standard',
          'Stimulus TLE Standard',
          'neurocaster',
          '{
            "processor": 2,
            "network": 2,
            "graphics": 2,
            "price": 700,
            "tags": []
          }'::jsonb
      ),
      (
          'stimulus_go',
          'Stimulus GO',
          'neurocaster',
          '{
            "processor": 2,
            "network": 2,
            "graphics": 1,
            "price": 600,
            "tags": [],
            "note": "Only gives -1 die to actions in the real world."
          }'::jsonb
      ),
      (
          'johnny_jolt_theme',
          'Johnny Jolt Theme',
          'neurocaster',
          '{
            "processor": 1,
            "network": 2,
            "graphics": 3,
            "price": 500,
            "tags": []
          }'::jsonb
      ),
      (
          'stimulus_tle_pro',
          'Stimulus TLE-PRO',
          'neurocaster',
          '{
            "processor": 3,
            "network": 3,
            "graphics": 3,
            "price": 1300,
            "tags": []
          }'::jsonb
      ),
      (
          'jury_rigged',
          'Jury-Rigged',
          'neurocaster',
          '{
            "processor": 1,
            "network": 1,
            "graphics": 1,
            "price": null,
            "tags": []
          }'::jsonb
      );

-- Drones

INSERT INTO tes.items (
    code,
    name,
    item_kind,
    data
) VALUES
      (
          'kids_drone_kid_kosmo',
          'Kids Drone "Kid Kosmo"',
          'drone',
          '{
            "strength": 3,
            "agility": 5,
            "hull": 8,
            "armor": 2,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "engaged",
            "price": 650,
            "tags": []
          }'::jsonb
      ),
      (
          'classic_gaming_drone_wally_wayne',
          'Classic Gaming Drone "Wally Wayne"',
          'drone',
          '{
            "strength": 4,
            "agility": 5,
            "hull": 9,
            "armor": 2,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 850,
            "tags": []
          }'::jsonb
      ),
      (
          'civilian_flyer_drone',
          'Civilian Flyer Drone',
          'drone',
          '{
            "strength": 3,
            "agility": 6,
            "hull": 9,
            "armor": 2,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 1200,
            "tags": ["capable_of_flight"]
          }'::jsonb
      ),
      (
          'battle_gaming_drone_johnny_jolt',
          'Battle Gaming Drone "Johnny Jolt"',
          'drone',
          '{
            "strength": 5,
            "agility": 4,
            "hull": 9,
            "armor": 3,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "medium",
            "price": 950,
            "tags": []
          }'::jsonb
      ),
      (
          'elite_trooper_gaming_drone',
          'Elite Trooper Gaming Drone',
          'drone',
          '{
            "strength": 5,
            "agility": 6,
            "hull": 11,
            "armor": 4,
            "damage_value": 2,
            "damage_kind": "direct",
            "range_min": "short",
            "range_max": "long",
            "price": 1150,
            "tags": []
          }'::jsonb
      ),
      (
          'jury_rigged_drone',
          'Jury-Rigged',
          'drone',
          '{
            "strength": 3,
            "agility": 3,
            "hull": 6,
            "armor": 1,
            "damage_value": 1,
            "damage_kind": "direct",
            "range_min": "engaged",
            "range_max": "short",
            "price": null,
            "tags": []
          }'::jsonb
      );

-- Gear

INSERT INTO tes.items (
    code,
    name,
    item_kind,
    data
) VALUES
      (
          'tools_general',
          'Tools, general',
          'gear',
          '{
            "bonus": 1,
            "price": 25,
            "tags": [],
            "note": "Can be used for any repairs (page 108)."
          }'::jsonb
      ),
      (
          'tools_vehicle',
          'Tools, vehicle',
          'gear',
          '{
            "bonus": 2,
            "price": 50,
            "tags": [],
            "note": "Can be used to repair vehicles."
          }'::jsonb
      ),
      (
          'tools_weapon',
          'Tools, weapon',
          'gear',
          '{
            "bonus": 2,
            "price": 100,
            "tags": [],
            "note": "Can be used to repair weapons."
          }'::jsonb
      ),
      (
          'tools_neurocaster',
          'Tools, neurocaster',
          'gear',
          '{
            "bonus": 2,
            "price": 50,
            "tags": [],
            "note": "Can be used to repair neurocasters."
          }'::jsonb
      ),
      (
          'first_aid_kit',
          'First aid kit',
          'gear',
          '{
            "bonus": 3,
            "price": 25,
            "tags": [],
            "uses": 5,
            "note": "Used to stabilize an Incapacitated person but requires the Medic talent."
          }'::jsonb
      ),
      (
          'surgical_instruments',
          'Surgical instruments',
          'gear',
          '{
            "bonus": 2,
            "price": 100,
            "tags": [],
            "note": "Gives gear bonus to performing surgery (page 84)."
          }'::jsonb
      ),
      (
          'pack_of_cigarettes',
          'Pack of cigarettes',
          'gear',
          '{
            "bonus": null,
            "price": 2,
            "tags": [],
            "uses": 4,
            "note": "Can be used once per Shift to recover 1 point of Hope, but also reduces Health by 1."
          }'::jsonb
      ),
      (
          'bottle_of_beer',
          'Bottle of beer',
          'gear',
          '{
            "bonus": null,
            "price": 2,
            "tags": [],
            "uses": 1,
            "note": "Can be used once per Day to recover 1 point of Hope."
          }'::jsonb
      ),
      (
          'bottle_of_hard_liquor',
          'Bottle of hard liquor',
          'gear',
          '{
            "bonus": null,
            "price": 5,
            "tags": [],
            "uses": 3,
            "note": "Can be used once per Shift to recover 1 point of Hope, but also reduces Health by 1."
          }'::jsonb
      ),
      (
          'pack_of_chewing_gum',
          'Pack of chewing gum',
          'gear',
          '{
            "bonus": 1,
            "price": 1,
            "tags": [],
            "uses": 3,
            "note": "Gives bonus to Empathy rolls when trying to be cool."
          }'::jsonb
      ),
      (
          'binoculars',
          'Binoculars',
          'gear',
          '{
            "bonus": 2,
            "price": 100,
            "tags": [],
            "note": "Used for Wits rolls to spot something at a distance."
          }'::jsonb
      ),
      (
          'neurine',
          'Neurine',
          'gear',
          '{
            "bonus": null,
            "price": 20,
            "tags": [],
            "uses": 1,
            "note": "Can be used once per Shift to recover 1 point of Hope. Roll for Wits after each use - if you fail, you become addicted and can only recover Hope in this way. Also known as dream glint."
          }'::jsonb
      ),
      (
          'food_canned',
          'Food, canned',
          'gear',
          '{
            "bonus": null,
            "price": 5,
            "tags": [],
            "uses": 1,
            "note": "Covers the daily food need for one person (page 89)."
          }'::jsonb
      ),
      (
          'clothes_outdoor',
          'Clothes, outdoor',
          'gear',
          '{
            "bonus": null,
            "price": 50,
            "tags": [],
            "note": "Keeps one person warm. If you do not have adequate clothes, you suffer the effects of cold (page 59)."
          }'::jsonb
      ),
      (
          'clothes_fine',
          'Clothes, fine',
          'gear',
          '{
            "bonus": 1,
            "price": 200,
            "tags": [],
            "note": "Gives bonus to Empathy rolls when trying to impress someone."
          }'::jsonb
      ),
      (
          'sleeping_bag',
          'Sleeping bag',
          'gear',
          '{
            "bonus": null,
            "price": 25,
            "tags": [],
            "note": "Allows one person to sleep comfortably outdoors, preventing sleep deprivation (page 89)."
          }'::jsonb
      ),
      (
          'shades',
          'Shades',
          'gear',
          '{
            "bonus": 1,
            "price": 20,
            "tags": [],
            "note": "Gives bonus to Empathy rolls when trying to be cool."
          }'::jsonb
      ),
      (
          'musical_instrument',
          'Musical instrument',
          'gear',
          '{
            "bonus": 2,
            "price": 100,
            "tags": [],
            "note": "Gives bonus to Empathy rolls but requires a Stretch of time and the Musician talent."
          }'::jsonb
      ),
      (
          'dog_pet',
          'Dog, pet',
          'gear',
          '{
            "bonus": null,
            "price": 100,
            "tags": [],
            "note": "Once per day, you can spend a Stretch with the dog to recover 1 Hope."
          }'::jsonb
      ),
      (
          'dog_guard',
          'Dog, guard',
          'gear',
          '{
            "bonus": null,
            "price": 250,
            "tags": [],
            "strength": 5,
            "agility": 4,
            "health": 9,
            "damage_value": 2,
            "damage_kind": "direct",
            "note": "Attacks on your command (action) and bites with base Damage 2."
          }'::jsonb
      ),
      (
          'book_fiction',
          'Book, fiction',
          'gear',
          '{
            "bonus": null,
            "price": 10,
            "tags": [],
            "note": "Once per day, you can spend a Stretch reading to recover 1 Hope."
          }'::jsonb
      ),
      (
          'book_religious',
          'Book, religious',
          'gear',
          '{
            "bonus": null,
            "price": 10,
            "tags": [],
            "note": "Once per day, you can spend a Stretch reading to recover 1 Hope."
          }'::jsonb
      ),
      (
          'book_medical',
          'Book, medical',
          'gear',
          '{
            "bonus": 1,
            "price": 30,
            "tags": [],
            "note": "Gives bonus to rolls for performing surgery (page 84)."
          }'::jsonb
      ),
      (
          'book_non_fiction',
          'Book, non-fiction',
          'gear',
          '{
            "bonus": 1,
            "price": 20,
            "tags": [],
            "note": "Gives bonus to Wits rolls if the subject is relevant."
          }'::jsonb
      ),
      (
          'newspaper',
          'Newspaper',
          'gear',
          '{
            "bonus": 1,
            "price": 0.5,
            "tags": [],
            "note": "Gives bonus to one Wits roll for anything related to current events."
          }'::jsonb
      ),
      (
          'walkman',
          'Walkman',
          'gear',
          '{
            "bonus": null,
            "price": 45,
            "tags": [],
            "note": "Once per day, you can spend a Stretch listening to recover 1 Hope."
          }'::jsonb
      ),
      (
          'spare_part',
          'Spare part',
          'gear',
          '{
            "bonus": null,
            "price": 100,
            "tags": [],
            "note": "Needed to repair an inoperable vehicle."
          }'::jsonb
      ),
      (
          'camera',
          'Camera',
          'gear',
          '{
            "bonus": null,
            "price": 200,
            "tags": [],
            "note": "Needs film."
          }'::jsonb
      ),
      (
          'pain_reliever',
          'Pain reliever',
          'gear',
          '{
            "bonus": null,
            "price": 3,
            "tags": [],
            "uses": 10,
            "note": "Heals 1 point of Health once per day, if not Incapacitated."
          }'::jsonb
      ),
      (
          'crowbar',
          'Crowbar',
          'gear',
          '{
            "bonus": 2,
            "price": 10,
            "tags": [],
            "damage_value": 1,
            "damage_kind": "direct",
            "note": "Gives bonus to Strength when breaking something. Can also be used as a weapon."
          }'::jsonb
      ),
      (
          'tent',
          'Tent',
          'gear',
          '{
            "bonus": null,
            "price": 75,
            "tags": [],
            "capacity": 4,
            "note": "Allows four people to sleep comfortably outdoors, preventing sleep deprivation (page 89)."
          }'::jsonb
      ),
      (
          'walkie_talkies',
          'Walkie-Talkies',
          'gear',
          '{
            "bonus": null,
            "price": 50,
            "tags": [],
            "note": "Allows communication up to about one mile."
          }'::jsonb
      ),
      (
          'gasoline_gallon',
          'Gasoline (gallon)',
          'gear',
          '{
            "bonus": null,
            "price": 1,
            "tags": [],
            "note": "A typical car runs about 20 miles per gallon."
          }'::jsonb
      ),
      (
          'jerrycan',
          'Jerrycan',
          'gear',
          '{
            "bonus": null,
            "price": 20,
            "tags": [],
            "capacity": 5,
            "capacity_unit": "gallon",
            "note": "Holds 5 gallons of gasoline."
          }'::jsonb
      ),
      (
          'vanadium_redox_battery',
          'Vanadium Redox battery',
          'gear',
          '{
            "bonus": null,
            "price": 50,
            "tags": ["not_commercial"],
            "note": "Powers drones, neurocasters and other electronic devices."
          }'::jsonb
      );