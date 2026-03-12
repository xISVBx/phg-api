package common

import (
	"context"
	"strings"

	"gorm.io/gorm"

	drepo "photogallery/api_go/internal/domain/repositories"
)

// NormalizePage normaliza los parámetros de paginación recibidos en QueryOptions.
//
// Reglas:
//   - Si Page <= 0, usa página 1.
//   - Si PageSize <= 0 o PageSize > 200, usa 25 como valor por defecto.
//
// Esto evita consultas inválidas o tamaños de página excesivos.
func NormalizePage(opts drepo.QueryOptions) (int, int) {
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 25
	}
	return page, pageSize
}

// ListWithQuery ejecuta una consulta paginada genérica con soporte para:
//
//   - búsqueda textual simple mediante opts.Q
//   - conteo total de registros filtrados
//   - ordenamiento dinámico
//   - limit/offset para paginación
//
// Parámetros:
//   - ctx: contexto de ejecución de la consulta.
//   - db: instancia base de GORM. Puede incluir filtros previos.
//   - model: modelo base usado en .Model(...).
//   - out: slice destino donde se cargarán los resultados.
//   - opts: opciones de consulta (paginación, búsqueda, ordenamiento).
//   - fields: campos permitidos para búsqueda textual con LIKE.
//   - allowedSorts: mapa de sorts permitidos (clave API -> expresión SQL segura).
//
// Comportamiento de búsqueda:
// Si opts.Q tiene valor y fields no está vacío, se construye una condición OR
// usando LOWER(campo) LIKE '%q%' sobre cada campo.
//
// Ejemplo:
// fields = ["code", "name"]
// q = "admin"
//
// Resultado equivalente:
// WHERE LOWER(code) LIKE '%admin%' OR LOWER(name) LIKE '%admin%'
//
// Comportamiento de ordenamiento:
//   - Por defecto ordena por "created_at_utc desc".
//   - Si opts.Sort viene informado, solo se aplica si existe en allowedSorts.
//   - opts.Dir solo permite "asc"; cualquier otro valor cae a "desc".
//
// Retorna:
//   - total: cantidad total de registros que cumplen el filtro, antes de paginar.
//   - error: error de GORM/SQL si ocurre.
func ListWithQuery[T any](
	ctx context.Context,
	db *gorm.DB,
	model any,
	out *[]T,
	opts drepo.QueryOptions,
	fields []string,
	allowedSorts map[string]string,
) (int64, error) {
	page, pageSize := NormalizePage(opts)
	tx := db.WithContext(ctx).Model(model)
	if q := strings.TrimSpace(opts.Q); q != "" && len(fields) > 0 {
		p := "%" + strings.ToLower(q) + "%"
		cond := ""
		args := []any{}
		for i, f := range fields {
			if i > 0 {
				cond += " OR "
			}
			cond += "LOWER(" + f + ") LIKE ?"
			args = append(args, p)
		}
		tx = tx.Where(cond, args...)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return 0, err
	}
	order := buildSafeOrder(opts, allowedSorts)
	if err := tx.Order(order).Limit(pageSize).Offset((page - 1) * pageSize).Find(out).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func buildSafeOrder(opts drepo.QueryOptions, allowedSorts map[string]string) string {
	const defaultOrder = "created_at_utc desc"

	if allowedSorts == nil {
		allowedSorts = map[string]string{}
	}

	sortKey := strings.TrimSpace(opts.Sort)
	if sortKey == "" {
		return defaultOrder
	}
	col, ok := allowedSorts[sortKey]
	if !ok || strings.TrimSpace(col) == "" {
		return defaultOrder
	}

	dir := strings.ToLower(strings.TrimSpace(opts.Dir))
	if dir != "asc" {
		dir = "desc"
	}
	return col + " " + dir
}
