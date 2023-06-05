workspace {

    model {
        !include https://raw.githubusercontent.com/ministryofjustice/opg-technical-guidance/main/dsl/poas/persons.dsl

        lpaCaseManagement = softwareSystem "LPA Case Management" "Existing System"

        lpaFrontend = softwareSystem "LPA Frontend" "Frontend UI for LPA case management" {
            container "ECS" "Provides UI for OPG users" "Go" "Component"
        }

        lpaFrontend -> lpaCaseManagement "Uses"

        caseworker -> lpaFrontend "Uses"
    }

    views {
        systemContext lpaFrontend "SystemContext" {
            include *
            autoLayout
        }

        container lpaFrontend {
            include *
            autoLayout
        }

        theme default

        styles {
            element "Existing System" {
                background #999999
                color #ffffff
            }
            element "Web Browser" {
                shape WebBrowser
            }
            element "Database" {
                shape Cylinder
            }
        }
    }
}
