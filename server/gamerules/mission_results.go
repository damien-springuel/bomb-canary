package gamerules

type missionResults map[mission]bool

func (m missionResults) copy() missionResults {
	c := make(missionResults)
	for k, v := range m {
		c[k] = v
	}
	return c
}

func (m missionResults) succeedMission(mission mission) missionResults {
	newResults := m.copy()
	newResults[mission] = true
	return newResults
}

func (m missionResults) failMission(mission mission) missionResults {
	newResults := m.copy()
	newResults[mission] = false
	return newResults
}

func (m missionResults) hasThreeSuccessesOrFailures() bool {
	successes, failures := 0, 0
	for _, success := range m {
		if success {
			successes += 1
		} else {
			failures += 1
		}
	}

	return successes == 3 || failures == 3
}
