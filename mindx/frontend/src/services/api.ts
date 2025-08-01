import axios from 'axios';
import type { RiskLevel, SortOption, Student } from '../types';

// Debug log
console.log('API module loaded');

const API_URL = 'http://localhost:8080';

// Add error handling for axios
axios.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export const evaluateStudents = async (): Promise<Student[]> => {
  try {
    console.log('Evaluating students...');
    const response = await axios.post(`${API_URL}/evaluate`);
    console.log('Evaluation response:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error in evaluateStudents:', error);
    throw error;
  }
};

export const getStudents = async (
  riskLevel?: RiskLevel,
  sortBy?: SortOption
): Promise<Student[]> => {
  try {
    let url = `${API_URL}/students`;
    const params = new URLSearchParams();
    
    if (riskLevel) {
      params.append('risk_level', riskLevel);
    }
    
    if (sortBy) {
      params.append('sort_by', sortBy);
    }
    
    const queryString = params.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
    
    console.log('Fetching students from URL:', url);
    const response = await axios.get(url);
    console.log('API response data:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error in getStudents:', error);
    throw error;
  }
};