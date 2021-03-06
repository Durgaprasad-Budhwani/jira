import React from 'react';
import { SimulatorInstaller, Integration, IProcessingDetail, IProcessingState, IInstalledLocation, ISelfManagedAgent, ISession } from '@pinpt/agent.websdk';
import IntegrationUI from './integration';

function App() {
	// check to see if we are running local and need to run in simulation mode
	if (window === window.parent && window.location.href.indexOf('localhost') > 0) {
		const integration: Integration = {
			name: 'Jira',
			description: 'The official Atlassian Jira integration for Pinpoint',
			tags: ['Issue Management'],
			installed: false,
			refType: 'jira',
			icon: 'https://pinpoint.com/images/integrations/Jira.svg',
			publisher: {
				name: 'Pinpoint',
				avatar: 'https://pinpoint.com/logo/logomark/blue.png',
				url: 'https://pinpoint.com'
			},
			uiURL: document.location.href,
		};

		const processingDetail: IProcessingDetail = {
			createdDate: Date.now() - (86400000 * 5) - 60000,
			processed: true,
			lastProcessedDate: Date.now() - (86400000 * 2),
			lastExportRequestedDate: Date.now() - ((86400000 * 5) + 60000),
			lastExportCompletedDate: Date.now() - (86400000 * 5),
			state: IProcessingState.IDLE,
			throttled: false,
			throttledUntilDate: Date.now() + 2520000,
			paused: false,
			location: IInstalledLocation.CLOUD
		};

		const selfManagedAgent: ISelfManagedAgent = {
			enrollment_id: '123',
			running: true,
		};

		const session: ISession = {
			customer: {
				id: '359d4a0ffac0329c',
				name: 'Pinpoint',
			},
			user: {
				id: '',
				name: 'Jeff Haynie',
				avatar_url: '',
			},
			env: 'edge',
			graphqlUrl: 'https://graph.api.edge.pinpoint.com/graphql',
			authUrl: 'https://auth.api.edge.pinpoint.com'
		};

		return <SimulatorInstaller id="f64c34f79f4b7994" integration={integration} processingDetail={processingDetail} selfManagedAgent={selfManagedAgent} session={session} />;
	}
	return <IntegrationUI />;
}

export default App;
