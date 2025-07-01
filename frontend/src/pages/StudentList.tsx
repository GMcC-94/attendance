// src/pages/StudentList.tsx
import { useEffect, useState } from 'react';
import axios from 'axios';

interface Student {
  id: number;
  name: string;
  beltGrade: string;
  dob: string;
}
/* eslint-disable  @typescript-eslint/no-explicit-any */
export default function StudentList() {
  const [students, setStudents] = useState<Student[]>([]);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState('');
  const [editBeltGrade, setEditBeltGrade] = useState('');
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

  const deleteStudent = async (id: number) => {
    if (!confirm("Are you sure you want to delete this student?")) return;

    try {
      await axios.delete(`api/v1/students/${id}`);
      setStudents(prev => prev.filter(student => student.id !== id));
    } catch (err) {
      console.error('Delete failed:', err)
      alert("Failed to delete student")
    }
  };

  const startEdit = (student: Student) => {
    setEditingId(student.id);
    setEditName(student.name);
    setEditBeltGrade(student.beltGrade);
  };

  const updateStudent = async () => {
    if (!editingId) return;

    const payload: any = {};
    if (editName.trim() !== '') payload.name = editName.trim();
    if (editBeltGrade.trim() !== '') payload.beltGrade = editBeltGrade.trim();

    if (Object.keys(payload).length === 0) {
      alert('Please provide at least one field to update.');
      return;
    }

    try {
      await axios.put(`/api/v1/students/${editingId}`, payload);

      setStudents(prev =>
        prev.map(s =>
          s.id === editingId
            ? { ...s, ...payload }
            : s
        )
      );

      setEditingId(null);
      setEditName('');
      setEditBeltGrade('');
    } catch (err) {
      console.error('Update failed:', err);
      alert('Failed to update student');
    }
  };

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4">Students</h2>
      <div className="flex gap-4 mb-4">
        <button
          className={`px-4 py-2 rounded ${filter === 'all' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setFilter('all')}
        >
          All
        </button>
        <button
          className={`px-4 py-2 rounded ${filter === 'adults' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setFilter('adults')}
        >
          Adults
        </button>
        <button
          className={`px-4 py-2 rounded ${filter === 'kids' ? 'bg-blue-600 text-white' : 'bg-gray-200'}`}
          onClick={() => setFilter('kids')}
        >
          Kids
        </button>
      </div>

      <table className="table-auto w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2 border">Name</th>
            <th className="p-2 border">Belt Grade</th>
            <th className="p-2 border">Date of Birth</th>
          </tr>
        </thead>
        <tbody>
          {students.map(student => (
            <tr key={student.id} className="text-center">
              {editingId === student.id ? (
                <>
                  <td className="border px-4 py-2">
                    <input
                      value={editName}
                      onChange={e => setEditName(e.target.value)}
                      className="border px-2 py-1"
                    />
                  </td>
                  <td className="border px-4 py-2">
                    <input
                      value={editBeltGrade}
                      onChange={e => setEditBeltGrade(e.target.value)}
                      className="border px-2 py-1"
                    />
                  </td>
                  <td className="border px-4 py-2">
                    <button
                      onClick={updateStudent}
                      className="bg-green-500 text-white px-2 py-1 mr-2 rounded"
                    >
                      Save
                    </button>
                    <button
                      onClick={() => setEditingId(null)}
                      className="bg-gray-500 text-white px-2 py-1 rounded"
                    >
                      Cancel
                    </button>
                  </td>
                </>
              ) : (
                <>
                  <td className="border px-4 py-2">{student.name}</td>
                  <td className="border px-4 py-2">{student.beltGrade}</td>
                  <td className="border px-4 py-2">{student.dob}</td>
                  <td className="border px-4 py-2">
                    <button
                      onClick={() => startEdit(student)}
                      className="bg-blue-500 text-white px-2 py-1 mr-2 rounded"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => deleteStudent(student.id)}
                      className="bg-red-500 text-white px-2 py-1 rounded"
                    >
                      Delete
                    </button>
                  </td>
                </>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
