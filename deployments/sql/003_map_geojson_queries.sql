-- =============================================================================
-- 基于 geojson (JSONB) 的地理位置相关查询示例
-- 依赖：map_features 表、002 种子数据（可选）
-- 说明：未使用 PostGIS；Point 用 coordinates[0/1]；线/面用坐标数组做范围或近似判断
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 辅助视图：仅 Point 类型，抽出经纬度便于范围查询与联表
-- -----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_map_point_coords AS
SELECT
  id,
  name,
  description,
  (geojson->'coordinates'->>0)::double precision AS lng,
  (geojson->'coordinates'->>1)::double precision AS lat,
  geojson
FROM map_features
WHERE geojson->>'type' = 'Point';

-- -----------------------------------------------------------------------------
-- 辅助视图：LineString，顶点展开为行（用于「顶点落在 bbox」等近似场景）
-- -----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_map_linestring_vertices AS
SELECT
  mf.id,
  mf.name,
  t.pos AS vertex_idx,
  (t.elem->>0)::double precision AS lng,
  (t.elem->>1)::double precision AS lat,
  mf.geojson
FROM map_features mf
CROSS JOIN LATERAL jsonb_array_elements(mf.geojson->'coordinates')
  WITH ORDINALITY AS t(elem, pos)
WHERE mf.geojson->>'type' = 'LineString';

-- -----------------------------------------------------------------------------
-- 场景 1：按 GeoJSON 几何类型筛选
-- -----------------------------------------------------------------------------
-- SELECT id, name, geojson->>'type' AS geom_type FROM map_features WHERE geojson->>'type' = 'Polygon';
-- SELECT id, name FROM map_features WHERE geojson->>'type' IN ('Point', 'LineString');

-- -----------------------------------------------------------------------------
-- 场景 2：JSON 包含 / 键存在（配合 GIN 索引 idx_map_features_geojson）
-- -----------------------------------------------------------------------------
-- 仅包含 type=Point 的文档结构匹配（示例）
-- SELECT id, name FROM map_features WHERE geojson @> '{"type":"Point"}';

-- -----------------------------------------------------------------------------
-- 场景 3：Point — 经纬度矩形范围（bounding box），例如故宫附近一块区域
-- -----------------------------------------------------------------------------
-- SELECT id, name, lng, lat
-- FROM v_map_point_coords
-- WHERE lng BETWEEN 116.396 AND 116.398
--   AND lat BETWEEN 39.915 AND 39.920;

-- 不依赖视图，直接对 JSONB：
-- SELECT id, name
-- FROM map_features
-- WHERE geojson->>'type' = 'Point'
--   AND (geojson->'coordinates'->>0)::double precision BETWEEN 116.396 AND 116.398
--   AND (geojson->'coordinates'->>1)::double precision BETWEEN 39.915 AND 39.920;

-- -----------------------------------------------------------------------------
-- 场景 4：Point — 与某固定点的「粗距离」排序（平面近似，非大地测量）
--      使用 (dx^2+dy^2) 排序，单位：度；仅适合小范围演示
-- -----------------------------------------------------------------------------
-- WITH ref AS (SELECT 116.397128::double precision AS lng, 39.916527::double precision AS lat)
-- SELECT m.id, m.name,
--        POWER((geojson->'coordinates'->>0)::double precision - ref.lng, 2)
--      + POWER((geojson->'coordinates'->>1)::double precision - ref.lat, 2) AS dist2_deg2
-- FROM map_features m, ref
-- WHERE m.geojson->>'type' = 'Point'
-- ORDER BY dist2_deg2
-- LIMIT 5;

-- -----------------------------------------------------------------------------
-- 场景 5：LineString — 是否存在顶点落在给定 bbox 内（路径与区域相交的简化判定）
-- -----------------------------------------------------------------------------
-- SELECT DISTINCT id, name
-- FROM v_map_linestring_vertices
-- WHERE lng BETWEEN 116.395 AND 116.399
--   AND lat BETWEEN 39.915 AND 39.918;

-- -----------------------------------------------------------------------------
-- 场景 6：Polygon — 外环任一点落在 bbox 内（多边形与矩形相交的宽松近似）
--      精确「点在多边形内」需 PostGIS ST_Contains 等
-- -----------------------------------------------------------------------------
-- SELECT DISTINCT mf.id, mf.name
-- FROM map_features mf
-- CROSS JOIN LATERAL jsonb_array_elements(mf.geojson->'coordinates'->0) AS pt
-- WHERE mf.geojson->>'type' = 'Polygon'
--   AND (pt->>0)::double precision BETWEEN 116.392 AND 116.401
--   AND (pt->>1)::double precision BETWEEN 39.914 AND 39.921;

-- -----------------------------------------------------------------------------
-- 场景 7：jsonpath（PostgreSQL 14+）— 按类型与坐标条件组合
-- -----------------------------------------------------------------------------
-- SELECT id, name
-- FROM map_features
-- WHERE jsonb_path_exists(
--   geojson,
--   '$ ? (@.type == "Point" && @.coordinates[0] > 116.39 && @.coordinates[0] < 116.40)'
-- );

-- -----------------------------------------------------------------------------
-- 可选：为 Point 的经度建立部分表达式索引，加速 bbox 查询（按需执行）
-- -----------------------------------------------------------------------------
-- CREATE INDEX IF NOT EXISTS idx_map_features_point_lng
--   ON map_features (((geojson->'coordinates'->>0)::double precision))
--   WHERE geojson->>'type' = 'Point';
