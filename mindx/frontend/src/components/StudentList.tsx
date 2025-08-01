import { useState, useEffect } from 'react';
import type { Student, RiskLevel, SortOption } from '../types';
import { getStudents, evaluateStudents } from '../services/api';
import StudentCard from './StudentCard';
import {
  Container,
  Typography,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Alert,
  CircularProgress
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';

// Debug log
console.log('StudentList component loaded');

const StudentList: React.FC = () => {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [riskFilter, setRiskFilter] = useState<RiskLevel>('');
  const [sortOption, setSortOption] = useState<SortOption>('');
  const [evaluating, setEvaluating] = useState<boolean>(false);

  const fetchStudents = async () => {
    console.log('Fetching students with filters:', { riskFilter, sortOption });
    setLoading(true);
    try {
      const data = await getStudents(riskFilter, sortOption);
      console.log('Students data received:', data);
      setStudents(data);
      setError(null);
    } catch (err) {
      console.error('Error fetching students:', err);
      setError('Failed to fetch students. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleEvaluate = async () => {
    console.log('Starting student evaluation');
    setEvaluating(true);
    try {
      const result = await evaluateStudents();
      console.log('Evaluation result:', result);
      await fetchStudents();
      setError(null);
    } catch (err) {
      console.error('Error evaluating students:', err);
      setError('Failed to evaluate students. Please try again.');
    } finally {
      setEvaluating(false);
    }
  };

  const handleRiskFilterChange = (event: SelectChangeEvent) => {
    setRiskFilter(event.target.value as RiskLevel);
  };

  const handleSortChange = (event: SelectChangeEvent) => {
    setSortOption(event.target.value as SortOption);
  };

  useEffect(() => {
    fetchStudents();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [riskFilter, sortOption]);

  return (
    <Container maxWidth="md">
      <Typography variant="h4" component="h1" gutterBottom sx={{ mt: 4 }}>
        Student Dropout Risk Dashboard
      </Typography>
      
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <FormControl variant="outlined" size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Risk Level</InputLabel>
            <Select
              value={riskFilter}
              onChange={handleRiskFilterChange}
              label="Risk Level"
            >
              <MenuItem value="">All Levels</MenuItem>
              <MenuItem value="HIGH">High Risk</MenuItem>
              <MenuItem value="MEDIUM">Medium Risk</MenuItem>
              <MenuItem value="LOW">Low Risk</MenuItem>
            </Select>
          </FormControl>
          
          <FormControl variant="outlined" size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Sort By</InputLabel>
            <Select
              value={sortOption}
              onChange={handleSortChange}
              label="Sort By"
            >
              <MenuItem value="">Default</MenuItem>
              <MenuItem value="risk_level">Risk Level (High → Low)</MenuItem>
              <MenuItem value="risk_level_asc">Risk Level (Low → High)</MenuItem>
              <MenuItem value="score">Score (High → Low)</MenuItem>
              <MenuItem value="score_asc">Score (Low → High)</MenuItem>
            </Select>
          </FormControl>
        </Box>
        
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleEvaluate}
          disabled={evaluating}
        >
          {evaluating ? <CircularProgress size={24} color="inherit" /> : 'Evaluate Students'}
        </Button>
      </Box>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
          <CircularProgress />
        </Box>
      ) : students.length > 0 ? (
        students.map((student) => (
          <StudentCard key={student.id} student={student} />
        ))
      ) : (
        <Alert severity="info">No students found.</Alert>
      )}
    </Container>
  );
};

export default StudentList;