package internal

import (
	"encoding/json"
	"os"

	varibale "github.com/yorukot/superfile/src/config"
)

// Pinned directory
func pinnedDirectory(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	unPinned := false

	jsonData, err := os.ReadFile(varibale.PinnedFilea)
	if err != nil {
		outPutLog("Pinned folder function read superfile data error", err)
	}

	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	if err != nil {
		outPutLog("Pinned folder function unmarshal superfile data error", err)
	}
	for i, other := range pinnedFolder {
		if other == panel.location {
			pinnedFolder = append(pinnedFolder[:i], pinnedFolder[i+1:]...)
			unPinned = true
		}
	}

	if !arrayContains(pinnedFolder, panel.location) && !unPinned {
		pinnedFolder = append(pinnedFolder, panel.location)
	}

	updatedData, err := json.Marshal(pinnedFolder)
	if err != nil {
		outPutLog("Pinned folder function updatedData superfile data error", err)
	}

	err = os.WriteFile(varibale.PinnedFilea, updatedData, 0644)
	if err != nil {
		outPutLog("Pinned folder function updatedData superfile data error", err)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

// Create new file panel
func createNewFilePanel(m model) model {
	if len(m.fileModel.filePanels) == m.fileModel.maxFilePanel {
		return m
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{
		location:        varibale.HomeDir,
		panelMode:       browserMode,
		focusType:       secondFocus,
		directoryRecord: make(map[string]directoryRecord),
		searchBar:       generateSearchBar(),
	})

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth) / Config.FilePreviewWidth
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.focusPanel)
	m.fileModel.width = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.filePanelFocusIndex++

	m.fileModel.maxFilePanel = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
	return m
}

// Close current focus file panel
func closeFilePanel(m model) model {
	if len(m.fileModel.filePanels) == 1 {
		return m
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {

			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth) / Config.FilePreviewWidth
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.width = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)

	m.fileModel.maxFilePanel = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
	return m
}

func toggleFilePreviewPanel(m model) model {
	m.fileModel.filePreview.open = !m.fileModel.filePreview.open
	m.fileModel.filePreview.width = 0
	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - sidebarWidth) / Config.FilePreviewWidth
		}
	}

	m.fileModel.width = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)

	m.fileModel.maxFilePanel = (m.fullWidth - sidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}

	return m
}

// Focus on next file panel
func nextFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

// Focus on previous file panel
func previousFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

// Focus on sidebar
func focusOnSideBar(m model) model {
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

// Focus on processbar
func focusOnProcessBar(m model) model {
	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

// focus on metadata
func focusOnMetadata(m model) model {
	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}
