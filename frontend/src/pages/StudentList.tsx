// src/pages/StudentList.tsx
import { useEffect, useState } from 'react';
import axios from 'axios';

interface Student {
  id: number;
  name: string;
  beltGrade: string;
  dob: string;
}

export default function StudentList() {
  const [students, setStudents] = useState<Student[]>([]);

  useEffect(() => {
    axios.get<Student[]>('/api/v1/students')
      .then((res) => setStudents(res.data))
      .catch((err) => console.error('Failed to fetch students', err));
  }, []);

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4">Students</h2>
      <table className="table-auto w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2 border">Name</th>
            <th className="p-2 border">Belt Grade</th>
            <th className="p-2 border">Date of Birth</th>
          </tr>
        </thead>
        <tbody>
          {students.map((student) => (
            <tr key={student.id}>
              <td className="p-2 border">{student.name}</td>
              <td className="p-2 border">{student.beltGrade}</td>
              <td className="p-2 border">{student.dob}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
