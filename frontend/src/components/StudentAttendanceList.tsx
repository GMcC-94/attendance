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
  const [filter, setFilter] = useState<'all' | 'adults' | 'kids'>('all');

  useEffect(() => {
    const fetchStudents = async () => {
      try {
        let endpoint = '/api/v1/students';
        if (filter === 'adults') endpoint = '/api/v1/students/adults';
        else if (filter === 'kids') endpoint = '/api/v1/students/kids';

        const res = await axios.get<Student[]>(endpoint);
        setStudents(res.data);
      } catch (err) {
        console.error('Failed to fetch students:', err);
      }
    };
    fetchStudents();
  }, [filter]);

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
     <div className="p-6 bg-gray-900 text-white rounded-lg shadow max-w-4xl mx-auto">
  <h2 className="text-2xl font-bold mb-4">Students</h2>
  <div className="flex gap-3 mb-6">
    <button
      className={`px-4 py-2 rounded transition ${
        filter === 'all'
          ? 'bg-blue-600 text-white hover:bg-blue-700'
          : 'bg-gray-700 text-gray-200 hover:bg-gray-600'
      }`}
      onClick={() => setFilter('all')}
    >
      All
    </button>
    <button
      className={`px-4 py-2 rounded transition ${
        filter === 'adults'
          ? 'bg-blue-600 text-white hover:bg-blue-700'
          : 'bg-gray-700 text-gray-200 hover:bg-gray-600'
      }`}
      onClick={() => setFilter('adults')}
    >
      Adults
    </button>
    <button
      className={`px-4 py-2 rounded transition ${
        filter === 'kids'
          ? 'bg-blue-600 text-white hover:bg-blue-700'
          : 'bg-gray-700 text-gray-200 hover:bg-gray-600'
      }`}
      onClick={() => setFilter('kids')}
    >
      Kids
    </button>
  </div>

  <table className="w-full border-collapse">
    <thead>
      <tr className="bg-gray-800">
        <th className="text-left py-3 px-4 font-semibold text-gray-200">Student Name</th>
        <th className="text-left py-3 px-4 font-semibold text-gray-200">Mark</th>
      </tr>
    </thead>
    <tbody>
      {students.map((student, idx) => (
        <tr
          key={student.id}
          className={`${idx % 2 === 0 ? 'bg-gray-800' : 'bg-gray-700'} border-b border-gray-600`}
        >
          <td className="py-3 px-4">{student.name}</td>
          <td className="py-3 px-4">
            <input
              type="checkbox"
              className="w-5 h-5 accent-blue-600"
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

    </div>

  );
};

export default StudentAttendanceList;
