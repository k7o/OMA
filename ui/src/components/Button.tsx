type ButtonProps = {
  text: string
  icon?: string
  onClick?: () => void
}

export const Button = ({ text, icon, onClick }: ButtonProps) => {
  return (
    <button
      class='px-4 py-2 mx-1 font-thin bg-blue-400 text-white rounded-md flex items-center'
      onClick={onClick}
    >
      {icon && <img src={icon} class='w-5 h-5' />}
      <span class='px-3'>{text}</span>
    </button>
  )
}
