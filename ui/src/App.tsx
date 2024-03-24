import { Router, Route, A, RouteSectionProps } from '@solidjs/router'
import { Playground } from './pages/Playground'
import { DecisionLogs } from './pages/DecisionLogs'
import { Home } from './pages/Home'
import { DataProvider } from './components/DataContext'

import HomeIcon from './assets/home.svg'
import PlayIcon from './assets/play-circle.svg'
import LogsIcon from './assets/logs.svg'

const App = () => {
  return (
    <DataProvider>
      <Router root={Page}>
        <Route path="/" component={Home} />
        <Route path="/play" component={Playground} />
        <Route path="/decision-logs" component={DecisionLogs} />
      </Router>
    </DataProvider>
  )
}

export default App

const Page = (props: RouteSectionProps<unknown>) => {
  return (
    <div class="flex w-screen h-screen">
      <div class="min-w-14 bg-[#eee] pt-16">
        <SidebarItem href="/" icon={HomeIcon} text="Home" />
        <SidebarItem href="/play" icon={PlayIcon} text="Play" />
        <SidebarItem href="/decision-logs" icon={LogsIcon} text="Logs" />
      </div>
      <div class="w-full">{props.children}</div>
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
    <A href={href} class="flex flex-col my-4 items-center border-2 text-sm font-thin">
      <img src={icon} alt={text} class="w-8 h-8 stroke-white color" />
      {text}
    </A>
  )
}
