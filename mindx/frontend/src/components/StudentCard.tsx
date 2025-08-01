import { useState } from 'react';
import type { Student } from '../types';
import {
  Card,
  CardContent,
  Typography,
  Chip,
  Box,
  Collapse,
  IconButton,
  Grid,
  List,
  ListItem,
  ListItemText,
  Divider
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';

// Debug log
console.log('StudentCard component loaded');

interface StudentCardProps {
  student: Student;
}

const getRiskLevelColor = (riskLevel: string | null): string => {
  switch (riskLevel) {
    case 'HIGH':
      return '#f44336'; // red
    case 'MEDIUM':
      return '#ff9800'; // orange
    case 'LOW':
      return '#4caf50'; // green
    default:
      return '#9e9e9e'; // gray
  }
};

const StudentCard: React.FC<StudentCardProps> = ({ student }) => {
  const [expanded, setExpanded] = useState(false);

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  const attendanceRate = student.attendance.length > 0
    ? (student.attendance.filter(a => a.status === 'ATTEND').length / student.attendance.length * 100).toFixed(1)
    : '100';

  const assignmentRate = student.assignments.length > 0
    ? (student.assignments.filter(a => a.submitted).length / student.assignments.length * 100).toFixed(1)
    : '100';

  const contactFailures = student.contacts.filter(c => c.status === 'FAILED').length;

  console.log('Rendering StudentCard for:', student.student_id);
  
  return (
    <Card sx={{ mb: 2, border: `1px solid ${getRiskLevelColor(student.dropout_risk_level)}` }}>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6" component="div">
            {student.student_name} ({student.student_id})
          </Typography>
          <Chip
            label={student.dropout_risk_level || 'Unknown'}
            sx={{
              bgcolor: getRiskLevelColor(student.dropout_risk_level),
              color: 'white',
              fontWeight: 'bold'
            }}
          />
        </Box>
        
        <Box mt={1}>
          <Typography variant="body2" color="text.secondary">
            Risk Score: {student.dropout_score !== null ? student.dropout_score : 'N/A'}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Note: {student.dropout_note || 'No notes'}
          </Typography>
        </Box>
        
        <Box mt={1} display="flex" justifyContent="space-between" alignItems="center">
          <Box>
            <Typography variant="caption">
              Attendance: {attendanceRate}% | Assignments: {assignmentRate}% | Contact Failures: {contactFailures}
            </Typography>
          </Box>
          <IconButton onClick={handleExpandClick}>
            {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
          </IconButton>
        </Box>
        
        <Collapse in={expanded} timeout="auto" unmountOnExit>
          <Box mt={2}>
            <Grid container spacing={2}>
              <Grid item xs={4}>
                <Typography variant="subtitle2">Attendance</Typography>
                <List dense>
                  {student.attendance.slice(0, 5).map((a, index) => (
                    <ListItem key={index} sx={{ py: 0 }}>
                      <ListItemText 
                        primary={a.date} 
                        secondary={a.status}
                        secondaryTypographyProps={{
                          color: a.status === 'ATTEND' ? 'success.main' : 'error.main'
                        }}
                      />
                    </ListItem>
                  ))}
                  {student.attendance.length > 5 && (
                    <ListItem sx={{ py: 0 }}>
                      <ListItemText primary={`+${student.attendance.length - 5} more...`} />
                    </ListItem>
                  )}
                </List>
              </Grid>
              
              <Divider orientation="vertical" flexItem />
              
              <Grid item xs={4}>
                <Typography variant="subtitle2">Assignments</Typography>
                <List dense>
                  {student.assignments.map((a, index) => (
                    <ListItem key={index} sx={{ py: 0 }}>
                      <ListItemText 
                        primary={a.name} 
                        secondary={a.submitted ? 'Submitted' : 'Not Submitted'}
                        secondaryTypographyProps={{
                          color: a.submitted ? 'success.main' : 'error.main'
                        }}
                      />
                    </ListItem>
                  ))}
                </List>
              </Grid>
              
              <Divider orientation="vertical" flexItem />
              
              <Grid item xs={3}>
                <Typography variant="subtitle2">Contact Attempts</Typography>
                <List dense>
                  {student.contacts.map((c, index) => (
                    <ListItem key={index} sx={{ py: 0 }}>
                      <ListItemText 
                        primary={c.date} 
                        secondary={c.status}
                        secondaryTypographyProps={{
                          color: c.status === 'SUCCESS' ? 'success.main' : 'error.main'
                        }}
                      />
                    </ListItem>
                  ))}
                  {student.contacts.length === 0 && (
                    <ListItem sx={{ py: 0 }}>
                      <ListItemText primary="No contact attempts" />
                    </ListItem>
                  )}
                </List>
              </Grid>
            </Grid>
          </Box>
        </Collapse>
      </CardContent>
    </Card>
  );
};

export default StudentCard;