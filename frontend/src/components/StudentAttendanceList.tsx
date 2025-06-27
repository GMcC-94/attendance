import React, { useEffect, useState } from 'react';
import axios from 'axios';

/* eslint-disable  @typescript-eslint/no-explicit-any */
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
        console.log("Fetched students:", res.data); 
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
      await axios.post(`/api/v1/students/${studentId}/attendance`, {
        attendedDays: [currentDay],
      });
      setMarked(prev => ({ ...prev, [studentId]: true }));
    } catch (error: any) {
      console.error('Error marking attendance:', error);
      alert(`Failed to mark attendance: ${error.response?.data || error.message}`);
    }
  };

  return (
    <div className="p-4 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Take Attendance for {currentDay}</h1>
      <table className="w-full border-collapse">
        <thead>
          <tr className="bg-gray-100">
            <th className="text-left py-2 px-4">Student Name</th>
            <th className="text-left py-2 px-4">Mark</th>
          </tr>
        </thead>
        <tbody>
          {students.map((student) => (
            <tr key={student.id} className="border-b">
              <td className="py-2 px-4">{student.name}</td>
              <td className="py-2 px-4">
                <input
                  type="checkbox"
                  disabled={marked[student.id]}
                  checked={marked[student.id] || false}
                  onChange={() => handleAttendance(student.id)}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default StudentAttendanceList;
