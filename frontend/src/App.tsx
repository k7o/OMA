import { DataProvider } from './components/DataContext'
import { Editor } from './components/Editor'
import { Header } from './components/Header'

const App = () => {


  return (
    <DataProvider>
      <Header />
      <Editor />
    </DataProvider>
  )
}

export default App
