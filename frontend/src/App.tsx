import Calculator from './components/Calculator'
import styles from './App.module.css'

function App() {
  return (
    <main className={styles.app}>
      <h1 className={styles.title}>Forge Abacus</h1>
      <Calculator />
    </main>
  )
}

export default App
