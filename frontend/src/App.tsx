import './App.css'
import StudentAttendanceList from './components/StudentAttendanceList'
import Navbar from './components/Navbar';

function App() {
  //const [count, setCount] = useState(0)

  return (
    <>
      <Navbar />
      <StudentAttendanceList />
    </>
  )
}

export default App
