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
            'Default Police Officer'
           'npc',
           '{
             "stats": {
               "strength": 4,
               "agility": 4,
               "wits": 3,
               "empathy": 3
             },
             "resources": {
               "health": 4,
               "hope": 2,
               "bliss": 0,
               "permanent_bliss": 0
             },
             "talents": ["Pistoleer"],
             "inventory": {
               "choices": [
                 {
                   "label": "Choose one weapon",
                   "choose": 1,
                   "item_type": "weapon",
                   "options": [
                     "handgun",
                     "shotgun"
                   ]
                }
              ]
             },
              "car": {
                "name": "Patrol car",
                "code": "4wd_car",
                "traits": ["powerful", "heavy"]
              },
            "cash": "$100"
           }'::jsonb
       );