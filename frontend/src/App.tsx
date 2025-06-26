import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import StudentAttendanceList from './components/StudentAttendanceList'

function App() {
  //const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-gray-100">
      <StudentAttendanceList />
    </div>
  )
}

export default App
