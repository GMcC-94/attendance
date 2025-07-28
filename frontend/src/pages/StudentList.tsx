import { useEffect, useState } from 'react';
import axios, { AxiosError } from 'axios';

interface Student {
  id: number;
  name: string;
  beltGrade: string;
  dob: string;
}

interface ApiError {
  error: string;
}

export default function StudentList() {
  const [students, setStudents] = useState<Student[]>([]);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState('');
  const [editBeltGrade, setEditBeltGrade] = useState('');
  const [filter, setFilter] = useState<'all' | 'adults' | 'kids'>('all');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStudents = async () => {
      try {
        let endpoint = '/api/v1/students';
        if (filter === 'adults') endpoint = '/api/v1/students/adults';
        else if (filter === 'kids') endpoint = '/api/v1/students/kids';

        const res = await axios.get<Student[]>(endpoint);
        setStudents(res.data);
      } catch (err) {
        const error = err as AxiosError<ApiError>;
        setError(error.response?.data?.error || 'Failed to fetch students');
        console.error('Failed to fetch students:', error);
      }
    };
    fetchStudents();
  }, [filter]);

  const deleteStudent = async (id: number) => {
    if (!confirm("Are you sure you want to delete this student?")) return;

    try {
      await axios.delete(`/api/v1/students/${id}`);
      setStudents(prev => prev.filter(student => student.id !== id));
    } catch (err) {
      const error = err as AxiosError<ApiError>;
      setError(error.response?.data?.error || 'Failed to delete student');
    }
  };

  const startEdit = (student: Student) => {
    setEditingId(student.id);
    setEditName(student.name);
    setEditBeltGrade(student.beltGrade);
  };

  const updateStudent = async () => {
    if (!editingId) return;

    const payload: Partial<Pick<Student, 'name' | 'beltGrade'>> = {};
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
      const error = err as AxiosError<ApiError>;
      setError(error.response?.data?.error || 'Failed to update student');
    }
  };

  return (
    <div className="p-6 bg-gray-900 text-white rounded-lg shadow max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">Students</h2>

      {error && <div className="bg-red-600 p-2 rounded mb-4">{error}</div>}

      <div className="flex gap-4 mb-4">
        {['all', 'adults', 'kids'].map((btn) => (
          <button
            key={btn}
            className={`px-4 py-2 rounded transition ${filter === btn
                ? 'bg-blue-600 text-white hover:bg-blue-700'
                : 'bg-gray-700 text-gray-200 hover:bg-gray-600'
              }`}
            onClick={() => setFilter(btn as 'all' | 'adults' | 'kids')}
          >
            {btn[0].toUpperCase() + btn.slice(1)}
          </button>
        ))}
      </div>

      <table className="table-auto w-full border border-gray-700">
        <thead>
          <tr className="bg-gray-800">
            <th className="p-3 border border-gray-700 text-left">Name</th>
            <th className="p-3 border border-gray-700 text-left">Belt Grade</th>
            <th className="p-3 border border-gray-700 text-left">Date of Birth</th>
            <th className="p-3 border border-gray-700 text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          {students.map(student => (
            <tr key={student.id} className="odd:bg-gray-800 even:bg-gray-700">
              {editingId === student.id ? (
                <>
                  <td className="border border-gray-700 px-3 py-2">
                    <input
                      value={editName}
                      onChange={e => setEditName(e.target.value)}
                      className="w-full bg-gray-800 border border-gray-600 rounded px-2 py-1 text-white focus:outline-none focus:border-blue-500"
                    />
                  </td>
                  <td className="border border-gray-700 px-3 py-2">
                    <input
                      value={editBeltGrade}
                      onChange={e => setEditBeltGrade(e.target.value)}
                      className="w-full bg-gray-800 border border-gray-600 rounded px-2 py-1 text-white focus:outline-none focus:border-blue-500"
                    />
                  </td>
                  <td className="border border-gray-700 px-3 py-2">{student.dob}</td>
                  <td className="border border-gray-700 px-3 py-2">
                    <button
                      onClick={updateStudent}
                      className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded mr-2"
                    >
                      Save
                    </button>
                    <button
                      onClick={() => setEditingId(null)}
                      className="bg-gray-600 hover:bg-gray-700 text-white px-3 py-1 rounded"
                    >
                      Cancel
                    </button>
                  </td>
                </>
              ) : (
                <>
                  <td className="border border-gray-700 px-3 py-2">{student.name}</td>
                  <td className="border border-gray-700 px-3 py-2">{student.beltGrade}</td>
                  <td className="border border-gray-700 px-3 py-2">{student.dob}</td>
                  <td className="border border-gray-700 px-3 py-2">
                    <button
                      onClick={() => startEdit(student)}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded mr-2"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => deleteStudent(student.id)}
                      className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded"
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
