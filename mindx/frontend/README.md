# Student Dropout Risk Dashboard Frontend

A React-based frontend application for visualizing and managing student dropout risk data.

## Features

- Interactive dashboard for monitoring student dropout risk
- Filtering and sorting capabilities for student data
- Material-UI components for a modern, responsive interface
- TypeScript for type safety and improved developer experience
- Integration with the Go backend API

## Tech Stack

- **React**: UI library for building component-based interfaces
- **TypeScript**: Type-safe JavaScript superset
- **Material-UI**: React component library implementing Google's Material Design
- **Axios**: HTTP client for API communication
- **Vite**: Next-generation frontend tooling for fast development and optimized builds

## Project Structure

```
frontend/
├── public/                # Static assets
├── src/
│   ├── assets/            # Images, fonts, and other assets
│   ├── components/        # React components
│   │   ├── StudentCard.tsx    # Individual student card component
│   │   └── StudentList.tsx    # List of students with filtering
│   ├── services/          # API services
│   │   └── api.ts         # API communication layer
│   ├── types/             # TypeScript type definitions
│   │   └── index.ts       # Shared type definitions
│   ├── App.tsx            # Main application component
│   ├── App.css            # Application styles
│   ├── main.tsx           # Application entry point
│   └── index.css          # Global styles
├── index.html             # HTML template
├── package.json           # Project dependencies and scripts
└── tsconfig.json          # TypeScript configuration
```

## Component Documentation

### StudentList Component

The `StudentList` component is the main container for displaying student data. It provides:

- Fetching and displaying student data from the API
- Filtering students by risk level (HIGH, MEDIUM, LOW)
- Sorting students by different criteria
- Triggering student risk evaluation

Key features:
- State management for loading, error, and data states
- Responsive layout for different screen sizes
- User-friendly error handling

### StudentCard Component

The `StudentCard` component displays individual student information including:

- Student personal information
- Risk level with color coding (HIGH: red, MEDIUM: orange, LOW: green)
- Detailed metrics for attendance, assignments, and communication
- Risk notes explaining the evaluation

## API Integration

The frontend communicates with the backend through a dedicated API service layer:

### API Service

The `api.ts` service provides:

1. `getStudents(riskLevel?, sortBy?)`: Fetches students with optional filtering and sorting
2. `evaluateStudents()`: Triggers risk evaluation on the backend

Features:
- Axios interceptors for global error handling
- Type-safe API responses using TypeScript interfaces
- Query parameter handling for filtering and sorting

## Type Definitions

The application uses TypeScript interfaces to ensure type safety:

```typescript
// Student data model
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

// Risk level type
export type RiskLevel = 'LOW' | 'MEDIUM' | 'HIGH' | '';

// Sort options
export type SortOption = 'risk_level' | 'risk_level_asc' | 'score' | 'score_asc' | '';
```

## Getting Started

### Prerequisites

- Node.js (v14 or later)
- npm or yarn

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd mindx/frontend
   ```
3. Install dependencies:
   ```bash
   npm install
   # or
   yarn
   ```

### Development

Start the development server:

```bash
npm run dev
# or
yarn dev
```

This will start the development server at http://localhost:3000 with hot module replacement (HMR) enabled.

### Building for Production

Build the application for production:

```bash
npm run build
# or
yarn build
```

This will generate optimized production files in the `dist` directory.

### Preview Production Build

Preview the production build locally:

```bash
npm run preview
# or
yarn preview
```

## Configuration

The frontend application can be configured through environment variables:

- `VITE_API_URL`: Backend API URL (default: http://localhost:8080)

Create a `.env` file in the frontend directory to override default settings:

```
VITE_API_URL=http://your-api-url
```

## ESLint Configuration

The project uses ESLint for code quality and consistency. The configuration can be expanded as described below:

```js
export default tseslint.config([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Add type-checking enabled rules
      ...tseslint.configs.recommendedTypeChecked,
      // Optionally use stricter rules
      ...tseslint.configs.strictTypeChecked,
      // Optionally add stylistic rules
      ...tseslint.configs.stylisticTypeChecked,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
])
```

## Best Practices

### Component Structure

- Keep components focused on a single responsibility
- Use TypeScript interfaces for component props
- Extract reusable logic to custom hooks
- Implement proper error handling and loading states

### State Management

- Use React hooks (useState, useEffect) for local component state
- Consider using Context API for shared state across components
- Implement proper data fetching patterns with loading and error states

### Styling

- Use Material-UI's styling system for consistent theming
- Leverage the theme provider for global style customization
- Use responsive design principles for mobile compatibility

### Performance Optimization

- Implement memoization for expensive calculations
- Use React.memo for components that render frequently
- Optimize re-renders by avoiding unnecessary state updates
- Implement virtualization for long lists of students