export interface AttendanceRecord {
  date: string;
  status: string;
}

export interface AssignmentRecord {
  date: string;
  name: string;
  submitted: boolean;
}

export interface ContactRecord {
  date: string;
  status: string;
}

export interface Student {
  id: string;
  student_id: string;
  student_name: string;
  attendance: AttendanceRecord[];
  assignments: AssignmentRecord[];
  contacts: ContactRecord[];
  dropout_score: number | null;
  dropout_risk_level: string | null;
  dropout_note: string | null;
  created_at: number;
  updated_at: number;
}

export type RiskLevel = 'LOW' | 'MEDIUM' | 'HIGH' | '';
export type SortOption = 'risk_level' | 'risk_level_asc' | 'score' | 'score_asc' | '';