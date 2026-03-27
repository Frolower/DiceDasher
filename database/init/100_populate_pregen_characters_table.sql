\connect character_db

INSERT INTO public.pregen_characters (
    system_name,
    name,
    description,
    type,
    template_data
)
VALUES (
           'tes',
           'Law Enforcement',
           'Default Police Officer',
           'npc',
           '{
             "name": "Law Enforcement",
             "archetype": "Police Officer",
             "favourite_song": "",
             "description": "Default Police Officer",
             "stats": {
               "strength": 4,
               "agility": 4,
               "wits": 3,
               "empathy": 3
             },
             "derivatives": {
               "health": 4,
               "hope": 2
             },
             "bliss": {
               "bliss": 0,
               "permanent": 0
             },
             "talents": ["Pistoleer"],
             "dream": "",
             "gear": [
               {
                 "name": "Handgun",
                 "code": "handgun",
                 "type": "weapon",
                 "bonus": 0,
                 "damage_value": 2,
                 "damage_kind": "ballistic",
                 "range_min": "0",
                 "range_max": "medium",
                 "price": 0,
                 "tags": []
               }
             ],
             "cash": 100,
             "tension": [],
             "conditions": {
               "injuries": [],
               "traumas": []
             },
             "vehicle": {
               "vehicle_type": "car",
               "model": "Patrol Car",
               "passenger": 4,
               "fuel": "gasoline",
               "description": "Standard police patrol car",
               "stats": {
                 "maneuverability": 2,
                 "speed": 3,
                 "hull": 4,
                 "armor": 1,
                 "traits": [],
                 "gear": []
               }
             }
           }'::jsonb
       );