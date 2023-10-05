import { A } from '@solidjs/router';

export function Home() {
  return (
    <ul>
      <li><A href='./ui'>Ui</A></li>
      <li><A href='./player'>Player</A></li>
    </ul>
  );
};

