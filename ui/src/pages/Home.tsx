import PlayIcon from '../assets/play-circle.svg'
import FileArchive from '../assets/file-archive.svg'
import GitFork from '../assets/git-fork.svg'

export const Home = () => {
  return (
    <div class="flex flex-col items-center justify-center w-full h-full">
      <h1 class="text-4xl text-gray-600 font-bold text mb-2">OMA - OPA Management Application</h1>
      <p class="text-lg mb-12 text-gray-600">
        This project aims to simplify the policy development workflow.
      </p>
      <div class="w-96 space-y-2">
        <ActionButton icon={PlayIcon} text="Create a playground" />
        <ActionButton icon={FileArchive} text="Create playground from bundle" />
        <ActionButton icon={GitFork} text="Create playground from git repository" />
      </div>
    </div>
  )
}

const ActionButton = (props: { text: string; icon: string; onClick?: () => void }) => {
  return (
    <button
      class="flex items-center w-full bg-gray-200 hover:bg-gray-400 text-gray-600 py-2 px-4 rounded"
      onClick={props.onClick}
    >
      <img src={props.icon} class="w-8 h-8 mr-4 stroke-white" />
      {props.text}
    </button>
  )
}
