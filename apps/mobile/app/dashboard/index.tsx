import { EmptyOrganizations } from '@/components/Overview/EmptyOrganizations';
import Overview from '@/components/Overview/Overview';
import { useOrganizationStore } from '@/stores/organizationStore';
import { useAppTheme } from '@/theme/theme';
import { View } from 'react-native';
import { ActivityIndicator } from 'react-native-paper';

export default function Index() {
    const theme = useAppTheme();
    const { organizations, initialized } = useOrganizationStore();

    if (!initialized) {
        return (
            <View
                style={{
                    flex: 1,
                    alignItems: 'center',
                    justifyContent: 'center',
                    backgroundColor: theme.colors.background,
                }}
            >
                <ActivityIndicator size="large" color={theme.colors.primary} />
            </View>
        );
    }

    if (organizations.length === 0) {
        return <EmptyOrganizations />;
    }

    return <Overview/>;
}
