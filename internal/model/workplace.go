package model

type WorkplaceResponse struct {
	Corporations []WorkplaceInfo `json:"hojin-infos"`
}

type WorkplaceInfo struct {
	CorporateNumber string      `json:"corporate_number"`
	Name            string      `json:"name"`
	Workplaces      []Workplace `json:"workplace_info"`
}

type Workplace struct {
	BaseMonth                     string `json:"base_month"`
	EmployeeNumber                string `json:"employee_number"`
	EmployeeNumberRegular         string `json:"employee_number_regular"`
	EmployeeNumberNonRegular      string `json:"employee_number_non_regular"`
	FemaleShareOfManager          string `json:"female_share_of_manager"`
	YearsOfService                string `json:"years_of_service"`
	AnnualSalary                  string `json:"annual_salary"`
	AverageContinuousServiceYears string `json:"average_continuous_service_years"`
	AverageAge                    string `json:"average_age"`
	MonthAverageOvertimeHours     string `json:"month_average_overtime_hours"`
	PaidHolidayUsageRate          string `json:"paid_holiday_usage_rate"`
}
