#!/bin/bash

# Guardian Insights Report
echo "📊 --- GUARDIAN ASSET INSIGHTS --- 📊"
echo ""

echo "🎯 Most Used Lenses:"
docker exec -it guardian-db-1 psql -U user -d guardian_db -t -c "
SELECT content->>'lens', COUNT(*) 
FROM daily_logs 
GROUP BY 1 ORDER BY 2 DESC;"

echo ""
echo "📸 Top Aperture Settings:"
docker exec -it guardian-db-1 psql -U user -d guardian_db -t -c "
SELECT content->>'aperture', COUNT(*) 
FROM daily_logs 
WHERE content->>'aperture' IS NOT NULL
GROUP BY 1 ORDER BY 2 DESC LIMIT 5;"