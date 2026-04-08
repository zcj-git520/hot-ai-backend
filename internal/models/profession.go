package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ProfessionCategory 职业分类模型
type ProfessionCategory struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uk_name" json:"name"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
	SortOrder   int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	Status      int       `gorm:"column:status;default:1;index:idx_status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ProfessionCategory) TableName() string {
	return "profession_categories"
}

// Profession 职业信息模型（核心表）
type Profession struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"column:name;type:varchar(100);not null;index:idx_name" json:"name"`
	Slug           string    `gorm:"column:slug;type:varchar(100);not null;uniqueIndex:uk_slug" json:"slug"`
	Icon           string    `gorm:"column:icon;type:varchar(20)" json:"icon,omitempty"`
	CategoryID     uint      `gorm:"column:category_id;index:idx_category" json:"category_id,omitempty"`
	Description    string    `gorm:"column:description;type:text" json:"description,omitempty"`
	RiskLevel      string    `gorm:"column:risk_level;type:enum('extreme','high','medium','low');default:'medium';index:idx_risk_level" json:"risk_level"`
	RiskScore      int       `gorm:"column:risk_score;default:50;index:idx_risk_score" json:"risk_score"`
	AutomationRate int       `gorm:"column:automation_rate;default:50" json:"automation_rate"`
	Status         int       `gorm:"column:status;default:1;index:idx_status" json:"status"`
	SortOrder      int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	PublishedAt    time.Time `gorm:"column:published_at" json:"published_at,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 关联数据（非数据库字段）
	CategoryName     string                    `gorm:"-" json:"category_name,omitempty"`
	ImpactAnalysis   *ProfessionImpactAnalysis `gorm:"-" json:"impact_analysis,omitempty"`
	TransitionAdvice *ProfessionTransitionAdvice `gorm:"-" json:"transition_advice,omitempty"`
	MarketData       *ProfessionMarketData     `gorm:"-" json:"market_data,omitempty"`
}

func (Profession) TableName() string {
	return "professions"
}

// ProfessionImpactAnalysis 职业影响分析模型
type ProfessionImpactAnalysis struct {
	ID             uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProfessionID   uint           `gorm:"column:profession_id;not null;uniqueIndex:uk_profession_id;index:idx_profession_id" json:"profession_id"`
	AffectedTasks  JSONStringList `gorm:"column:affected_tasks;type:json" json:"affected_tasks,omitempty"`
	SafeTasks      JSONStringList `gorm:"column:safe_tasks;type:json" json:"safe_tasks,omitempty"`
	SafeSkills     JSONStringList `gorm:"column:safe_skills;type:json" json:"safe_skills,omitempty"`
	ImpactTimeline JSONStringMap  `gorm:"column:impact_timeline;type:json" json:"impact_timeline,omitempty"`
	ImpactSummary  string         `gorm:"column:impact_summary;type:text" json:"impact_summary,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// JSONStringList 自定义类型: 数据库存JSON字符串, 序列化时输出数组
type JSONStringList []string

func (j *JSONStringList) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (j JSONStringList) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j JSONStringList) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(j))
}

// JSONStringMap 自定义类型: 数据库存JSON对象字符串, 序列化时输出对象
type JSONStringMap map[string]string

func (j *JSONStringMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (j JSONStringMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j JSONStringMap) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal(map[string]string(j))
}

func (ProfessionImpactAnalysis) TableName() string {
	return "profession_impact_analysis"
}

// ProfessionTransitionAdvice 职业转型建议模型
type ProfessionTransitionAdvice struct {
	ID                uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProfessionID      uint           `gorm:"column:profession_id;not null;uniqueIndex:uk_profession_id;index:idx_profession_id" json:"profession_id"`
	TransitionPaths   JSONStringList `gorm:"column:transition_paths;type:json" json:"transition_paths,omitempty"`
	RecommendedSkills JSONStringList `gorm:"column:recommended_skills;type:json" json:"recommended_skills,omitempty"`
	RecommendedTools  JSONStringList `gorm:"column:recommended_tools;type:json" json:"recommended_tools,omitempty"`
	RecommendedPaths  JSONPathList   `gorm:"column:recommended_paths;type:json" json:"recommended_paths,omitempty"`
	RelatedArticles   JSONIntArray   `gorm:"column:related_articles;type:json" json:"related_articles,omitempty"`
	AdviceSummary     string         `gorm:"column:advice_summary;type:text" json:"advice_summary,omitempty"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// JSONPathList 推荐路径列表
type JSONPathList []map[string]interface{}

func (j *JSONPathList) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (j JSONPathList) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j JSONPathList) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]map[string]interface{}(j))
}

// JSONIntArray JSON整数数组
type JSONIntArray []int

func (j *JSONIntArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (j JSONIntArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j JSONIntArray) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]int(j))
}

func (ProfessionTransitionAdvice) TableName() string {
	return "profession_transition_advice"
}

// ProfessionMarketData 职业市场数据模型
type ProfessionMarketData struct {
	ID                     uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProfessionID           uint      `gorm:"column:profession_id;not null;uniqueIndex:uk_profession_id;index:idx_profession_id" json:"profession_id"`
	MarketTrend            string    `gorm:"column:market_trend;type:enum('growing','stable','declining');default:'stable';index:idx_market_trend" json:"market_trend"`
	MarketTrendDescription string    `gorm:"column:market_trend_description;type:text" json:"market_trend_description,omitempty"`
	SalaryImpact           string    `gorm:"column:salary_impact;type:enum('positive','neutral','negative');default:'neutral'" json:"salary_impact"`
	SalaryChangeRate       float64   `gorm:"column:salary_change_rate;type:decimal(5,2)" json:"salary_change_rate,omitempty"`
	AvgSalary              float64   `gorm:"column:avg_salary;type:decimal(10,2)" json:"avg_salary,omitempty"`
	JobDemandTrend         string    `gorm:"column:job_demand_trend;type:varchar(50)" json:"job_demand_trend,omitempty"`
	SupplyDemandRatio      float64   `gorm:"column:supply_demand_ratio;type:decimal(3,2)" json:"supply_demand_ratio,omitempty"`
	DataSource             string    `gorm:"column:data_source;type:varchar(255)" json:"data_source,omitempty"`
	DataUpdateDate         time.Time `gorm:"column:data_update_date;type:date" json:"data_update_date,omitempty"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ProfessionMarketData) TableName() string {
	return "profession_market_data"
}

// RiskLevelInfo 风险等级展示信息（静态数据）
type RiskLevelInfo struct {
	ID          string `json:"id"`
	Level       string `json:"level"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Color       string `json:"color"`
	MinScore    int    `json:"min_score"`
	MaxScore    int    `json:"max_score"`
}
