import './App.css'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import "tailwindcss";
import './index.css'
import StudentList from './pages/StudentList';
import Layout from './components/Layout';
import StudentAttendanceList from './components/StudentAttendanceList';
import CreateStudentForm from './components/CreateStudentForm';
import UploadLogo from './pages/UploadLogo';

function App() {

return (
   <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="attendance" element={<StudentAttendanceList />} />
          <Route path="students" element={<StudentList />} />
          <Route path="students/create" element={<CreateStudentForm />} />
           <Route path="upload-logo" element={<UploadLogo />} />
        </Route>
      </Routes>
    </Router>
  );
};

export default App;

