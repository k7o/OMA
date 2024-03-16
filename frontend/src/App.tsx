import { Router, Route, Navigate, A } from '@solidjs/router'
import { Home as Playground } from './pages/Playground'
import { DecisionLogs } from './pages/DecisionLogs'

import PlayIcon from './assets/play.png'
import LogsIcon from './assets/logs.svg'
import { JSX } from 'solid-js'

const App = () => {
  return (
    <Router>
      <Route path='/' component={() => <Navigate href='/play' />} />
      <Route path='/play' component={Page(Playground)} />
      <Route path='/decision-logs' component={Page(DecisionLogs)} />
    </Router>
  )
}

export default App

const Page = (children: () => JSX.Element) => {
  return () => (
    <div class='flex flex-row w-full h-full'>
      <div class='w-14 bg-[#eee] pt-16'>
        <SidebarItem href='/play' icon={PlayIcon} text='Play' />
        <SidebarItem href='/decision-logs' icon={LogsIcon} text='Logs' />
      </div>
      <div class='w-screen h-screen'>{children()}</div>
    </div>
  )
}

type SidebarItemProps = {
  href: string
  icon: string
  text: string
}

const SidebarItem = ({ href, icon, text }: SidebarItemProps) => {
  return (
    <A
      href={href}
      class='flex flex-col my-4 items-center border-2 text-sm font-thin'
    >
      <img src={icon} alt={text} class='w-8 h-8 stroke-white' />
      {text}
    </A>
  )
}
