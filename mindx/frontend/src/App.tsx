import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import StudentList from './components/StudentList';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#2196f3',
    },
    secondary: {
      main: '#f50057',
    },
  },
});

function App() {
  console.log("App component rendered");
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <StudentList />
    </ThemeProvider>
  );
}

export default App;