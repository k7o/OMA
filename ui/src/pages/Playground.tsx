import { Editor } from '../components/Editor'
import { Header } from '../components/Header'

export const Playground = () => {
  console.log(import.meta.env)

  return (
    <>
      <Header />
      <Editor />
    </>
  )
}
