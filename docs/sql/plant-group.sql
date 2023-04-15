-- Select plant group by ID
SELECT PG.ID AS PLANT_GROUP_ID,
    PG.NAME AS PLANT_GROUP_NAME,
    PG.DESCRIPTION AS PLANT_GROUP_DESCRIPTION
FROM PLANT_GROUP PG
WHERE PG.ID = ?;

-- Select all plant group IDs
SELECT ID
FROM PLANT_GROUP;
