export const SideBar = () => {
  return (
    <div class='flex flex-col h-screen w-1/6 bg-gray-900 text-white'>
      <div class='flex items-center justify-center h-16 bg-gray-800'>
        <h1 class='text-2xl'>Sidebar</h1>
      </div>
      <div class='flex flex-col items-center justify-center h-full'>
        <div class='flex items-center justify-center h-16 w-full'>
          <a href='/home'>Home</a>
        </div>
        <div class='flex items-center justify-center h-16 w-full'>
          <a href='/decision-logs'>Decision Logs</a>
        </div>
      </div>
    </div>
  )
}
