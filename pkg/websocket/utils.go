package websocket

func Find(a []*Client, u *Client) int {
	for i, n := range a {
		if u.Connection == n.Connection {
			return i
		}
	}
	return -1
}

func RemoveIndex(s []*Client, index int) []*Client {
	return append(s[:index], s[index+1:]...) //check how this works
}

func Contains(a []*Client, x *Client) (bool, int) {
	for i, u := range a {
		if x.Username == u.Username {
			return true, i
		}
	}
	return false, 0
}