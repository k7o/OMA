import { A } from '@solidjs/router'
import { DataProvider } from '../components/DataContext'
import { Editor } from '../components/Editor'
import { Header } from '../components/Header'
import { SideBar } from '../components/Sidebar'

export const Home = () => {
  return (
    <DataProvider>
      <Header />
      <Editor />
    </DataProvider>
  )
}
