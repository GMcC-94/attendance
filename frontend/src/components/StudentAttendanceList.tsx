import React, { useEffect, useState } from 'react';
import axios from 'axios';

interface Student {
  id: number;
  name: string;
  beltGrade: string;
  dob: string;
}

const allowedDays = new Set(['Monday', 'Tuesday', 'Wednesday', 'Friday']);
const currentDay = new Date().toLocaleDateString('en-US', { weekday: 'long' });

const StudentAttendanceList: React.FC = () => {
  const [students, setStudents] = useState<Student[]>([]);
  const [marked, setMarked] = useState<Record<number, boolean>>({});

  useEffect(() => {
    axios.get('/api/v1/students')
      .then(res => {
        console.log("Fetched students:", res.data); // ✅ ADD THIS LINE
        if (Array.isArray(res.data)) {
          setStudents(res.data);
        } else {
          console.error("Expected an array but got:", res.data);
        }
      })
      .catch(err => console.error('Failed to fetch students', err));
  }, []);

  const handleAttendance = async (studentId: number) => {
    if (!allowedDays.has(currentDay)) {
      alert(`Attendance cannot be taken on ${currentDay}`);
      return;
    }

    try {
      await axios.post(`/app/v1/students/${studentId}/attendance`, {
        attendedDays: [currentDay],
      });
      setMarked(prev => ({ ...prev, [studentId]: true }));
    } catch (error: any) {
      console.error('Error marking attendance:', error);
      alert(`Failed to mark attendance: ${error.response?.data || error.message}`);
    }
  };

 return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6 text-center">Take Attendance for {currentDay}</h1>
      <table className="min-w-full table-auto border-collapse border border-gray-300 shadow-sm rounded-md">
        <thead>
          <tr className="bg-gray-100">
            <th className="border border-gray-300 px-6 py-3 text-left">Name</th>
            <th className="border border-gray-300 px-6 py-3 text-center">Present</th>
          </tr>
        </thead>
        <tbody>
          {students.map(student => (
            <tr key={student.id} className="even:bg-gray-50 hover:bg-gray-100">
              <td className="border border-gray-300 px-6 py-3">{student.name}</td>
              <td className="border border-gray-300 px-6 py-3 text-center">
                <input
                  type="checkbox"
                  checked={!!marked[student.id]}
                  disabled={!!marked[student.id]}
                  onChange={() => handleAttendance(student.id)}
                  className="cursor-pointer w-5 h-5"
                />
              </td>
            </tr>
          ))}
          {students.length === 0 && (
            <tr>
              <td colSpan={2} className="text-center py-6 text-gray-500">
                No students found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default StudentAttendanceList;
