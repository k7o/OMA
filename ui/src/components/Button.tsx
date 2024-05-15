import { Component, JSX } from 'solid-js'

export const Button = (props: {
  text: string
  icon?: Component<JSX.SvgSVGAttributes<SVGSVGElement>>
  onClick?: () => void
}) => {
  return (
    <button
      class="px-4 py-2 mx-1 font-thin bg-blue-400 text-white rounded-md flex items-center"
      onClick={props.onClick}
    >
      {props.icon && <props.icon class="w-5 h-5 text-white" />}
      <span class="px-3">{props.text}</span>
    </button>
  )
}
