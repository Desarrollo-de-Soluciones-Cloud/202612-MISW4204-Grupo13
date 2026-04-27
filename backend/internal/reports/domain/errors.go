package domain

import "errors"

var (
	ErrReportInvalidInput         = errors.New("entrada invalida")
	ErrReportNotFound             = errors.New("reporte no encontrado")
	ErrReportWorkspaceIDRequired  = errors.New("workspace_id es requerido")
	ErrReportWeekIDRequired       = errors.New("week_id es requerido")
	ErrReportAssignmentIDRequired = errors.New("assignment_id es requerido")
	ErrReportUserIDRequired       = errors.New("user_id es requerido")
	ErrReportSummaryRequired      = errors.New("summary es requerido")
	ErrReportFilePathRequired     = errors.New("file_path es requerido")
	ErrReportWeekNotFound         = errors.New("semana no encontrada")
	ErrReportWorkspaceNotFound    = errors.New("workspace no encontrado")
	ErrReportFileNotFound         = errors.New("archivo del reporte no encontrado")
	ErrReportForbidden            = errors.New("no tiene permisos para acceder al reporte")
)
