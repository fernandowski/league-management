import {useLeagueData} from "@/hooks/useLeagueData";
import {useOrganizationStore} from "@/stores/organizationStore";
import {StyleSheet, Text, View} from "react-native";
import {Select} from "@/components/Select/Select";
import {useEffect} from "react";


export interface Props {
    onChange: (value: string) => void
    selected?: string
}

export default function LeagueDropdown(props: Props) {
    const {organization} = useOrganizationStore();
    const {fetchData, data} = useLeagueData()

    useEffect(() => {
        if (organization !== null) {
            fetchData({organization_id: organization, limit: 0, offset: 0, term: ""})
        }
    }, [organization])

    const onSelectChange = (value: string) => {
        props.onChange(value)
    }

    let selected = data.length > 0 ? data[0].id : null
    if (props.selected) {
        selected = props.selected
    }

    return (
        <View style={[styles.select]}>
            <Text>Leagues: </Text>
            <Select
                onChange={onSelectChange}
                data={data.map(league => ({label: league.name, value: league.id}))}
                selected={selected}
            />
        </View>
    )
}


const styles = StyleSheet.create({
    select: {
        flex: 1,
        flexDirection: "row",
        marginTop: 18
    }
});
