import { A, useLocation } from '@solidjs/router'

import PlayIcon from '../assets/play-circle.svg'
import FileArchive from '../assets/file-archive.svg'
import GitFork from '../assets/git-fork.svg'
import { Component, JSX } from 'solid-js'

export const Home = () => {
  return (
    <div class="flex flex-col items-center justify-center w-full h-full">
      <h1 class="text-4xl text-gray-600 font-bold text mb-2 max-w-2xl text-center">
        OMA - OPA Management Application
      </h1>
      <p class="text-lg mb-12 mx-6 text-center text-gray-600">
        This project aims to simplify the policy development workflow.
      </p>
      <div class="w-96 space-y-2">
        <ActionButton href="/new/empty" icon={PlayIcon} text="Create a playground" />
        <ActionButton
          href="/new/from-bundle"
          icon={FileArchive}
          text="Create a playground from a bundle"
        />
        <ActionButton
          href="/new/from-repository"
          icon={GitFork}
          text="Create playground from git repository"
        />
      </div>
    </div>
  )
}

export const ActionButton = (props: {
  text: string
  icon: Component<JSX.SvgSVGAttributes<SVGSVGElement>>
  href: string
}) => {
  return (
    <A
      href={props.href}
      class="flex items-center w-full bg-gray-200 hover:bg-gray-300 text-gray-600 py-2 px-4 rounded"
    >
      <props.icon class="w-8 h-8 mr-4" />
      {props.text}
    </A>
  )
}
