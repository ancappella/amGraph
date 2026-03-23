-- 示例地图数据：与 map_features.geojson (JSONB) 兼容，符合 GeoJSON RFC 7946
-- 坐标示例为 WGS84（经度, 纬度），位于北京城区一带，便于联调展示
--
-- 已存在同名要素时跳过整批插入，便于重复执行脚本而不产生重复行

INSERT INTO map_features (name, description, geojson)
SELECT name, description, geojson::jsonb
FROM (VALUES
  (
    '天安门广场',
    '地标点：示例 Point',
    '{"type":"Point","coordinates":[116.397128,39.916527]}'
  ),
  (
    '故宫博物院',
    '地标点：示例 Point',
    '{"type":"Point","coordinates":[116.397026,39.918058]}'
  ),
  (
    '示例步行路径',
    '折线路径：LineString，连接两处地标',
    '{"type":"LineString","coordinates":[[116.3958,39.9152],[116.3965,39.9158],[116.3971,39.9164],[116.3978,39.9171]]}'
  ),
  (
    '示例矩形区域',
    '闭合多边形：Polygon（首尾坐标需相同）',
    '{"type":"Polygon","coordinates":[[[116.393,39.915],[116.400,39.915],[116.400,39.920],[116.393,39.920],[116.393,39.915]]]}'
  ),
  (
    '景山公园',
    '小型绿地范围：Polygon',
    '{"type":"Polygon","coordinates":[[[116.3910,39.9280],[116.3955,39.9280],[116.3955,39.9310],[116.3910,39.9310],[116.3910,39.9280]]]}'
  ),
  (
    '北海公园入口',
    '补充点要素',
    '{"type":"Point","coordinates":[116.3889,39.9252]}'
  )
) AS v(name, description, geojson)
WHERE NOT EXISTS (
  SELECT 1 FROM map_features WHERE name = '天安门广场'
);
